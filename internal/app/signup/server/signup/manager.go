/*
 * Copyright 2020 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package signup

import (
	"context"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-organization-manager-go"
	"github.com/nalej/grpc-signup-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

const DefaultStorageAllocationSize = 100 * 1024 * 1024
const DefaultStorageAllocationSizeDesc = "Default Storage Size (bytes)"

// DefaultRoles defines the map of roles that will be automatically created by the system by default.
var DefaultRoles = map[string][]grpc_authx_go.AccessPrimitive{
	"Owner": {
		grpc_authx_go.AccessPrimitive_ORG,
	},
	"Operator": {
		grpc_authx_go.AccessPrimitive_PROFILE,
		grpc_authx_go.AccessPrimitive_RESOURCES,
	},
	"Developer": {
		grpc_authx_go.AccessPrimitive_PROFILE,
		grpc_authx_go.AccessPrimitive_APPS,
	},
	"AppCluster": {
		grpc_authx_go.AccessPrimitive_APPCLUSTEROPS,
	},
	"NalejAdmin": {
		grpc_authx_go.AccessPrimitive_ORG,
		grpc_authx_go.AccessPrimitive_ORG_MNGT,
		grpc_authx_go.AccessPrimitive_RESOURCES_MNGT,
	},
}

// InternalRoles contains the relationship of which roles are application managed (no human involved).
var InternalRoles = map[string]bool{
	"Owner":      false,
	"Operator":   false,
	"Developer":  false,
	"NalejAdmin": false,
	"AppCluster": true,
}

// Manager structure with the required providers for cluster operations.
type Manager struct {
	OrgClient     grpc_organization_manager_go.OrganizationsClient
	UserClient    grpc_user_manager_go.UserManagerClient
	ClusterClient grpc_infrastructure_go.ClustersClient
	AppClient     grpc_application_go.ApplicationsClient
}

// NewManager creates a Manager using a set of providers.
func NewManager(
	orgClient grpc_organization_manager_go.OrganizationsClient,
	userClient grpc_user_manager_go.UserManagerClient,
	clusterClient grpc_infrastructure_go.ClustersClient,
	appClient grpc_application_go.ApplicationsClient,
) Manager {
	return Manager{orgClient, userClient, clusterClient, appClient}
}

func (m *Manager) SignupOrganization(signupRequest *grpc_signup_go.SignupOrganizationRequest) (*grpc_organization_manager_go.Organization, error) {

	addOrganizationRequest := &grpc_organization_go.AddOrganizationRequest{
		Name:        signupRequest.OrganizationName,
		Email:       signupRequest.OrganizationEmail,
		FullAddress: signupRequest.OrganizationFullAddress,
		City:        signupRequest.OrganizationCity,
		State:       signupRequest.OrganizationState,
		Country:     signupRequest.OrganizationCountry,
		ZipCode:     signupRequest.OrganizationZipCode,
		PhotoBase64: signupRequest.OrganizationPhotoBase64,
	}
	orgCreated, err := m.OrgClient.AddOrganization(context.Background(), addOrganizationRequest)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating organization")
		return nil, err
	}
	log.Debug().Str("organizationID", orgCreated.OrganizationId).Msg("Organization has been created")

	// create organization settings
	_, err = m.OrgClient.AddSetting(context.Background(), &grpc_organization_go.AddSettingRequest{
		OrganizationId: orgCreated.OrganizationId,
		Key:            grpc_organization_go.AllowedSettingKey_DEFAULT_STORAGE_SIZE.String(),
		Value:          fmt.Sprintf("%d", DefaultStorageAllocationSize),
		Description:    DefaultStorageAllocationSizeDesc,
	})
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating settings")
	} else {
		log.Debug().Str("organizationID", orgCreated.OrganizationId).Str("setting", grpc_organization_go.AllowedSettingKey_DEFAULT_STORAGE_SIZE.String()).Msg("Setting added")
	}

	ownerRoleID, nalejAdminRoleID, err := m.createRoles(orgCreated.OrganizationId)
	if err != nil {
		// TODO Rollback required
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating roles")
		return nil, err
	}

	// Add Nalej administrator
	addNalejAdminRequest := &grpc_user_manager_go.AddUserRequest{
		OrganizationId: orgCreated.OrganizationId,
		Email:          signupRequest.NalejadminEmail,
		Password:       signupRequest.NalejadminPassword,
		Name:           signupRequest.NalejadminName,
		RoleId:         *nalejAdminRoleID,
	}
	err = addUser(err, m, addNalejAdminRequest, orgCreated)
	if err != nil {
		return nil, err
	}

	// Add owner
	addOwnerRequest := &grpc_user_manager_go.AddUserRequest{
		OrganizationId: orgCreated.OrganizationId,
		Email:          signupRequest.OwnerEmail,
		Password:       signupRequest.OwnerPassword,
		Name:           signupRequest.OwnerName,
		RoleId:         *ownerRoleID,
	}
	err = addUser(err, m, addOwnerRequest, orgCreated)
	if err != nil {
		return nil, err
	}
	return orgCreated, nil
}

func (m *Manager) createRoles(organizationID string) (*string, *string, error) {
	var ownerRoleID string
	var nalejAdminRoleID string
	for name, primitives := range DefaultRoles {
		internal, found := InternalRoles[name]
		if !found {
			return nil, nil, derrors.NewInternalError("cannot determine if role is internal")
		}
		addRoleRequest := &grpc_user_manager_go.AddRoleRequest{
			OrganizationId: organizationID,
			Name:           name,
			Description:    "Auto generate role",
			Internal:       internal,
			Primitives:     primitives,
		}
		added, err := m.UserClient.AddRole(context.Background(), addRoleRequest)
		if err != nil {
			return nil, nil, err
		}
		switch name {
		case "Owner":
			ownerRoleID = added.RoleId
		case "NalejAdmin":
			nalejAdminRoleID = added.RoleId
		}
		log.Debug().Str("organizationID", organizationID).Str("roleID", added.RoleId).Msg("Role has been created")
	}
	return &ownerRoleID, &nalejAdminRoleID, nil
}

// ListOrganizations returns the list of organizations in the system.
func (m *Manager) ListOrganizations(request *grpc_signup_go.SignupInfoRequest) (*grpc_signup_go.OrganizationsList, error) {
	orgs, err := m.OrgClient.ListOrganizations(context.Background(), &grpc_common_go.Empty{})
	if err != nil {
		return nil, err
	}
	result := make([]*grpc_signup_go.OrganizationInfo, 0, len(orgs.Organizations))
	for _, org := range orgs.Organizations {
		info, err := m.extendOrganizationInfo(org)
		if err != nil {
			return nil, err
		}
		result = append(result, info)
	}
	return &grpc_signup_go.OrganizationsList{
		Organizations: result,
	}, err
}

func (m *Manager) extendOrganizationInfo(org *grpc_organization_manager_go.Organization) (*grpc_signup_go.OrganizationInfo, error) {
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: org.OrganizationId,
	}

	clusters, err := m.ClusterClient.ListClusters(context.Background(), orgID)
	if err != nil {
		return nil, err
	}
	descriptors, err := m.AppClient.ListAppDescriptors(context.Background(), orgID)
	if err != nil {
		return nil, err
	}
	instances, err := m.AppClient.ListAppInstances(context.Background(), orgID)
	if err != nil {
		return nil, err
	}
	return &grpc_signup_go.OrganizationInfo{
		OrganizationId:    org.OrganizationId,
		Name:              org.Name,
		Created:           org.Created,
		NumberUsers:       org.NumUsers,
		NumberClusters:    int32(len(clusters.Clusters)),
		NumberDescriptors: int32(len(descriptors.Descriptors)),
		NumberInstances:   int32(len(instances.Instances)),
	}, nil
}

// GetOrganizationInfo retrieves the information about an organization.
func (m *Manager) GetOrganizationInfo(organizationID *grpc_organization_go.OrganizationId) (*grpc_signup_go.OrganizationInfo, error) {
	org, err := m.OrgClient.GetOrganization(context.Background(), organizationID)
	if err != nil {
		return nil, err
	}
	return m.extendOrganizationInfo(org)
}

// DeleteOrganization removes an organization from the system.
func (m *Manager) RemoveOrganization(organizationID *grpc_organization_go.OrganizationId) error {
	log.Info().Str("organizationID", organizationID.OrganizationId).Msg("Removing organization")
	// Undeploy running apps
	// Delete descriptors
	// Delete nodes
	// Delete clusters
	// Delete users
	// Delete roles
	// Delete organization
	panic("implement me")
}

func addUser(err error, m *Manager, addNalejAdminRequest *grpc_user_manager_go.AddUserRequest, orgCreated *grpc_organization_manager_go.Organization) error {
	nalejAdminAdded, err := m.UserClient.AddUser(context.Background(), addNalejAdminRequest)
	if err != nil {
		log.Error().Str("roleID", addNalejAdminRequest.RoleId).Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating user")
		return err
	}
	log.Debug().Str("organizationID", orgCreated.OrganizationId).Str("role", nalejAdminAdded.RoleName).Msg("User has been created")
	return nil
}

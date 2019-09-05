/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package signup

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-signup-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

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
		grpc_authx_go.AccessPrimitive_ORGMNGT,
		grpc_authx_go.AccessPrimitive_RESOURCESMNGT,
	},
}

// InternalRoles contains the relationship of which roles are internal.
var InternalRoles = map[string]bool{
	"Owner":      false,
	"Operator":   false,
	"Developer":  false,
	"AppCluster": true,
	"NalejAdmin": true,
}

// Manager structure with the required providers for cluster operations.
type Manager struct {
	OrgClient     grpc_organization_go.OrganizationsClient
	UserClient    grpc_user_manager_go.UserManagerClient
	ClusterClient grpc_infrastructure_go.ClustersClient
	AppClient     grpc_application_go.ApplicationsClient
}

// NewManager creates a Manager using a set of providers.
func NewManager(
	orgClient grpc_organization_go.OrganizationsClient,
	userClient grpc_user_manager_go.UserManagerClient,
	clusterClient grpc_infrastructure_go.ClustersClient,
	appClient grpc_application_go.ApplicationsClient,
) Manager {
	return Manager{orgClient, userClient, clusterClient, appClient}
}

func (m *Manager) SignupOrganization(signupRequest *grpc_signup_go.SignupOrganizationRequest) (*grpc_organization_go.Organization, error) {

	addOrganizationRequest := &grpc_organization_go.AddOrganizationRequest{
		Name: signupRequest.OrganizationName,
	}
	orgCreated, err := m.OrgClient.AddOrganization(context.Background(), addOrganizationRequest)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating organization")
		return nil, err
	}
	log.Debug().Str("organizationID", orgCreated.OrganizationId).Msg("Organization has been created")

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
	nalejAdminAdded, err := m.UserClient.AddUser(context.Background(), addNalejAdminRequest)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating Nalej admin")
		return nil, err
	}
	log.Debug().Str("organizationID", orgCreated.OrganizationId).Str("role", nalejAdminAdded.RoleName).Msg("User has been created")

	// Add owner
	addOwnerRequest := &grpc_user_manager_go.AddUserRequest{
		OrganizationId: orgCreated.OrganizationId,
		Email:          signupRequest.OwnerEmail,
		Password:       signupRequest.OwnerPassword,
		Name:           signupRequest.OwnerName,
		RoleId:         *ownerRoleID,
	}
	ownerAdded, err := m.UserClient.AddUser(context.Background(), addOwnerRequest)
	if err != nil {
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating owner")
		return nil, err
	}
	log.Debug().Str("organizationID", orgCreated.OrganizationId).Str("role", ownerAdded.RoleName).Msg("User has been created")
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

func (m *Manager) extendOrganizationInfo(org *grpc_organization_go.Organization) (*grpc_signup_go.OrganizationInfo, error) {
	orgID := &grpc_organization_go.OrganizationId{
		OrganizationId: org.OrganizationId,
	}
	users, err := m.UserClient.ListUsers(context.Background(), orgID)
	if err != nil {
		return nil, err
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
		NumberUsers:       int32(len(users.Users)),
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

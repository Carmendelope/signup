/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package signup

import (
	"context"
	"github.com/nalej/grpc-authx-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-signup-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

var DefaultRoles = map[string][]grpc_authx_go.AccessPrimitive{
	"Owner": []grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_ORG},
	"Developer" : []grpc_authx_go.AccessPrimitive{grpc_authx_go.AccessPrimitive_RESOURCES, grpc_authx_go.AccessPrimitive_APPS},
}

// Manager structure with the required providers for cluster operations.
type Manager struct {
	OrgClient grpc_organization_go.OrganizationsClient
	UserClient grpc_user_manager_go.UserManagerClient
}

// NewManager creates a Manager using a set of providers.
func NewManager(orgClient grpc_organization_go.OrganizationsClient, userClient grpc_user_manager_go.UserManagerClient) Manager {
	return Manager{orgClient, userClient}
}

func (m * Manager) SignupOrganization(signupRequest *grpc_signup_go.SignupOrganizationRequest) (*grpc_organization_go.Organization, error) {

	addOrganizationRequest := &grpc_organization_go.AddOrganizationRequest{
		Name:                 signupRequest.OrganizationName,
	}
	orgCreated, err := m.OrgClient.AddOrganization(context.Background(), addOrganizationRequest)
	if err != nil{
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating organization")
		return nil, err
	}
	log.Debug().Str("organizationID", orgCreated.OrganizationId).Msg("Organization has been created")

	ownerRoleID, err := m.createRoles(orgCreated.OrganizationId)
	if err != nil{
		// TODO Rollback required
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating roles")
		return nil, err
	}

	addUserRequest := &grpc_user_manager_go.AddUserRequest{
		OrganizationId:       orgCreated.OrganizationId,
		Email:                signupRequest.OwnerEmail,
		Password:             signupRequest.OwnerPassword,
		Name:                 signupRequest.OwnerName,
		RoleId:               *ownerRoleID,
	}
	userAdded, err := m.UserClient.AddUser(context.Background(), addUserRequest)
	if err != nil{
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error creating user")
		return nil, err
	}
	log.Debug().Str("organizationID", orgCreated.OrganizationId).Str("rol", userAdded.RoleName).Msg("User has been created")
	return orgCreated, nil
}

func (m * Manager) createRoles(organizationID string) (*string, error) {
	ownerRoleID := ""
	for name, primitives := range DefaultRoles{
		addRoleRequest := &grpc_user_manager_go.AddRoleRequest{
			OrganizationId:       organizationID,
			Name:                 name,
			Description:          "Auto generate rol",
			Primitives:           primitives,
		}
		added, err := m.UserClient.AddRole(context.Background(), addRoleRequest)
		if err != nil{
			return nil, err
		}
		if name == "Owner" {
			ownerRoleID = added.RoleId
		}
		log.Debug().Str("organizationID", organizationID).Str("roleID", added.RoleId).Msg("Rol has been created")
	}
	return &ownerRoleID, nil
}
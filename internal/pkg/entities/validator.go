/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-signup-go"
)

func ValidOrganizationId(organizationID *grpc_organization_go.OrganizationId) derrors.Error {
	if organizationID.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("organization_id must be provided")
	}
	return nil
}

func ValidSignupOrganizationRequest(signupRequest *grpc_signup_go.SignupOrganizationRequest) derrors.Error {
	if signupRequest.OrganizationName == "" {
		return derrors.NewInvalidArgumentError("organization_name must be provided")
	}
	if signupRequest.OwnerEmail == "" {
		return derrors.NewInvalidArgumentError("owner_email must be provided")
	}
	if signupRequest.OwnerName == "" {
		return derrors.NewInvalidArgumentError("owner_name must be provided")
	}
	if signupRequest.OwnerPassword == "" {
		return derrors.NewInvalidArgumentError("owner_password must be provided")
	}
	if signupRequest.NalejadminEmail == "" {
		return derrors.NewInvalidArgumentError("nalejadmin_email must be provided")
	}
	if signupRequest.NalejadminName == "" {
		return derrors.NewInvalidArgumentError("nalejadmin_name must be provided")
	}
	if signupRequest.NalejadminPassword == "" {
		return derrors.NewInvalidArgumentError("nalejadmin_password must be provided")
	}
	return nil
}

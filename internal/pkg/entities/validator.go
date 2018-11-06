package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-signup-go"
)

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
	return nil
}

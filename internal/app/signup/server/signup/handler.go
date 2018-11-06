/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package signup

import (
	"context"
	"github.com/nalej/grpc-signup-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/signup/internal/pkg/entities"
)

// Handler structure for the cluster requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler{
	return &Handler{manager}
}

func (h*Handler) SignupOrganization(ctx context.Context, signupRequest *grpc_signup_go.SignupOrganizationRequest) (*grpc_signup_go.SignupOrganizationResponse, error) {
	vErr := entities.ValidSignupOrganizationRequest(signupRequest)
	if vErr != nil{
		return nil, conversions.ToGRPCError(vErr)
	}
	organization, err := h.Manager.SignupOrganization(signupRequest)
	if err != nil{
		return nil, err
	}
	return &grpc_signup_go.SignupOrganizationResponse{
		OrganizationId:       organization.OrganizationId,
	}, nil
}
/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package signup

import (
	"context"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"

	"github.com/nalej/grpc-signup-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/signup/internal/pkg/entities"
)

// Handler structure for the cluster requests.
type Handler struct {
	Manager Manager
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// SignupOrganization register a new organization in the system with a new
// user as the owner.
func (h *Handler) SignupOrganization(ctx context.Context, signupRequest *grpc_signup_go.SignupOrganizationRequest) (*grpc_signup_go.SignupOrganizationResponse, error) {
	vErr := entities.ValidSignupOrganizationRequest(signupRequest)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	organization, err := h.Manager.SignupOrganization(signupRequest)
	if err != nil {
		return nil, err
	}
	return &grpc_signup_go.SignupOrganizationResponse{
		OrganizationId: organization.OrganizationId,
	}, nil
}

// ListOrganizations returns the list of organizations in the system.
func (h *Handler) ListOrganizations(ctx context.Context, _ *grpc_common_go.Empty) (*grpc_signup_go.OrganizationsList, error){
	return h.Manager.ListOrganizations()
}

// GetOrganizationInfo retrieves the information about an organization.
func (h *Handler) GetOrganizationInfo(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_signup_go.OrganizationInfo, error){
	vErr := entities.ValidOrganizationId(organizationID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.GetOrganizationInfo(organizationID)
}
// DeleteOrganization removes an organization from the system.
func (h *Handler) RemoveOrganization(ctx context.Context, organizationID *grpc_organization_go.OrganizationId) (*grpc_common_go.Success, error){
	vErr := entities.ValidOrganizationId(organizationID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.RemoveOrganization(organizationID)
	if err != nil{
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}


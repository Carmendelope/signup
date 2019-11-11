/*
 * Copyright 2019 Nalej
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
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-organization-go"
	"github.com/rs/zerolog/log"

	"github.com/nalej/grpc-signup-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/signup/internal/pkg/entities"
)

// Handler structure for the cluster requests.
type Handler struct {
	Manager              Manager
	CheckPresharedSecret bool
	PresharedSecret      string
}

// NewHandler creates a new Handler with a linked manager.
func NewHandler(manager Manager, checkPresharedSecret bool, presharedSecret string) *Handler {
	return &Handler{manager, checkPresharedSecret, presharedSecret}
}

func (h *Handler) checkPresharedSecret(found string) derrors.Error {
	if !h.CheckPresharedSecret {
		return nil
	}
	if h.PresharedSecret != found {
		return derrors.NewPermissionDeniedError("invalid preshared secret")
	}
	return nil
}

// SignupOrganization register a new organization in the system with a new
// user as the owner.
func (h *Handler) SignupOrganization(ctx context.Context, signupRequest *grpc_signup_go.SignupOrganizationRequest) (*grpc_signup_go.SignupOrganizationResponse, error) {
	sErr := h.checkPresharedSecret(signupRequest.PresharedSecret)
	if sErr != nil {
		log.Error().Str("trace", conversions.ToDerror(sErr).DebugReport()).Msg("error validating secret")
		return nil, sErr
	}
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
func (h *Handler) ListOrganizations(ctx context.Context, request *grpc_signup_go.SignupInfoRequest) (*grpc_signup_go.OrganizationsList, error) {
	sErr := h.checkPresharedSecret(request.PresharedSecret)
	if sErr != nil {
		log.Error().Str("trace", conversions.ToDerror(sErr).DebugReport()).Msg("error validating secret")
		return nil, sErr
	}
	return h.Manager.ListOrganizations(request)
}

// GetOrganizationInfo retrieves the information about an organization.
func (h *Handler) GetOrganizationInfo(ctx context.Context, request *grpc_signup_go.SignupInfoRequest) (*grpc_signup_go.OrganizationInfo, error) {
	sErr := h.checkPresharedSecret(request.PresharedSecret)
	if sErr != nil {
		log.Error().Str("trace", conversions.ToDerror(sErr).DebugReport()).Msg("error validating secret")
		return nil, sErr
	}
	organizationID := &grpc_organization_go.OrganizationId{
		OrganizationId: request.OrganizationId,
	}
	vErr := entities.ValidOrganizationId(organizationID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	return h.Manager.GetOrganizationInfo(organizationID)
}

// DeleteOrganization removes an organization from the system.
func (h *Handler) RemoveOrganization(ctx context.Context, request *grpc_signup_go.SignupInfoRequest) (*grpc_common_go.Success, error) {
	sErr := h.checkPresharedSecret(request.PresharedSecret)
	if sErr != nil {
		log.Error().Str("trace", conversions.ToDerror(sErr).DebugReport()).Msg("error validating secret")
		return nil, sErr
	}
	organizationID := &grpc_organization_go.OrganizationId{
		OrganizationId: request.OrganizationId,
	}
	vErr := entities.ValidOrganizationId(organizationID)
	if vErr != nil {
		return nil, conversions.ToGRPCError(vErr)
	}
	err := h.Manager.RemoveOrganization(organizationID)
	if err != nil {
		return nil, err
	}
	return &grpc_common_go.Success{}, nil
}

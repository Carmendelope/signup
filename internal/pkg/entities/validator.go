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
	if signupRequest.OrganizationFullAddress == "" {
		return derrors.NewInvalidArgumentError("organization_full_address must be provided")
	}
	if signupRequest.OrganizationCity == "" {
		return derrors.NewInvalidArgumentError("organization_city must be provided")
	}
	if signupRequest.OrganizationState == "" {
		return derrors.NewInvalidArgumentError("organization_state must be provided")
	}
	if signupRequest.OrganizationCountry == "" {
		return derrors.NewInvalidArgumentError("organization_country must be provided")
	}
	if signupRequest.OrganizationZipCode == "" {
		return derrors.NewInvalidArgumentError("organization_zip_code must be provided")
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

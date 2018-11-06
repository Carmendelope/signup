/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-signup-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type SignupCli struct {
	client grpc_signup_go.SignupClient
}

func NewSignupCli(signupAddress string) (* SignupCli, derrors.Error) {
	sConn, err := grpc.Dial(signupAddress, grpc.WithInsecure())
	if err != nil{
		return nil, derrors.AsError(err, "cannot create connection with the system model")
	}
	c := grpc_signup_go.NewSignupClient(sConn)
	return &SignupCli{c}, nil
}

func (s * SignupCli) SignupOrganization(orgName string, ownerEmail string, ownerName string, ownerPassword string){
	signupRequest := &grpc_signup_go.SignupOrganizationRequest{
		OrganizationName:     orgName,
		OwnerEmail:           ownerEmail,
		OwnerName:            ownerName,
		OwnerPassword:        ownerPassword,
	}
	response, err := s.client.SignupOrganization(context.Background(), signupRequest)
	if err != nil {
		dErr := conversions.ToDerror(err)
		log.Error().Str("err", dErr.Error()).Msg("cannot signup organization")
		log.Error().Str("trace", conversions.ToDerror(err).DebugReport()).Msg("error")
		return
	}
	log.Info().Str("organizationID", response.OrganizationId).Msg("organization has been added")
}
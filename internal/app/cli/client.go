/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-signup-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

//SignupCli with necessary data to create a new client
type SignupCli struct {
	client grpc_signup_go.SignupClient
	PresharedSecret string
}

//NewSignupCli connects to the Signup service send signup requests
func NewSignupCli(signupAddress string, caPath string, clientCertPath string, clientKeyPath string, presharedSecret string) (*SignupCli, derrors.Error) {
	var sConn *grpc.ClientConn
	var dErr error
	if caPath == "" || clientCertPath == "" || clientKeyPath == "" {
		log.Warn().Msg("Using client without certificates")
		sConn, dErr = grpc.Dial(signupAddress, grpc.WithInsecure())
		if dErr != nil {
			return nil, derrors.AsError(dErr, "cannot create connection with the signup service")
		}
	} else {
		creds, err := getTLSConfig(caPath, clientCertPath, clientKeyPath)
		if err != nil {
			return nil, err
		}
		sConn, dErr = grpc.Dial(signupAddress, grpc.WithTransportCredentials(creds))
		if dErr != nil {
			return nil, derrors.AsError(dErr, "cannot create connection with the signup service")
		}
	}

	c := grpc_signup_go.NewSignupClient(sConn)
	return &SignupCli{c, presharedSecret}, nil
}

//SignupOrganization sends the request to create a new Organization based on the arguments given
func (s *SignupCli) SignupOrganization(orgName string, ownerEmail string, ownerName string, ownerPassword string) {
	signupRequest := &grpc_signup_go.SignupOrganizationRequest{
		OrganizationName: orgName,
		OwnerEmail:       ownerEmail,
		OwnerName:        ownerName,
		OwnerPassword:    ownerPassword,
		PresharedSecret: s.PresharedSecret,
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

func getTLSConfig(caPath string, clientCertPath string, clientKeyPath string) (credentials.TransportCredentials, derrors.Error) {
	rootCAs := x509.NewCertPool()

	if caPath != "" {
		caCert, err := ioutil.ReadFile(caPath)
		if err != nil {
			return nil, derrors.NewInternalError("Error loading CA certificate")
		}
		rootCAs.AppendCertsFromPEM(caCert)
	}

	clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, derrors.NewInternalError("Error loading client certificate")
	}

	tlsConfig := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{clientCert},
	}

	return credentials.NewTLS(tlsConfig), nil
}

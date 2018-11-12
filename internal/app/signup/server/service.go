/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package server

import (
	"fmt"
	"net"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-organization-go"
	"github.com/nalej/grpc-signup-go"
	"github.com/nalej/grpc-user-manager-go"
	"github.com/nalej/grpc-utils/pkg/tools"
	"github.com/nalej/signup/internal/app/signup/server/signup"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//Service struct to define the gRPC Server and its configuration
type Service struct {
	Configuration Config
	Server        *tools.GenericGRPCServer
}

// NewService creates a new system model service.
func NewService(conf Config) *Service {
	return &Service{
		conf,
		tools.NewGenericGRPCServer(uint32(conf.Port)),
	}
}

//Clients definition
type Clients struct {
	orgClient  grpc_organization_go.OrganizationsClient
	userClient grpc_user_manager_go.UserManagerClient
}

//GetClients gets a new instance of Clients with an active client of every type defined
func (s *Service) GetClients() (*Clients, derrors.Error) {
	smConn, err := grpc.Dial(s.Configuration.SystemModelAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the system model")
	}

	uConn, err := grpc.Dial(s.Configuration.UserManagerAddress, grpc.WithInsecure())
	if err != nil {
		return nil, derrors.AsError(err, "cannot create connection with the user manager")
	}

	oClient := grpc_organization_go.NewOrganizationsClient(smConn)
	uClient := grpc_user_manager_go.NewUserManagerClient(uConn)
	log.Debug().Str("smConn", smConn.GetState().String()).Str("uConn", uConn.GetState().String()).Msg("connections have been created")

	return &Clients{oClient, uClient}, nil
}

// Run the service, launch the REST service handler.
func (s *Service) Run() error {
	vErr := s.Configuration.Validate()
	if vErr != nil {
		log.Fatal().Str("err", vErr.DebugReport()).Msg("invalid configuration")
	}

	s.Configuration.Print()

	return s.LaunchGRPC()
}

//LaunchGRPC creates the gRPC server, register the necessary handlers and serves it
func (s *Service) LaunchGRPC() error {
	clients, cErr := s.GetClients()
	if cErr != nil {
		log.Fatal().Str("err", cErr.DebugReport()).Msg("cannot generate clients")
		return cErr
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Configuration.Port))
	if err != nil {
		log.Fatal().Errs("failed to listen: %v", []error{err})
	}

	manager := signup.NewManager(clients.orgClient, clients.userClient)
	handler := signup.NewHandler(manager)

	grpcServer := grpc.NewServer()
	grpc_signup_go.RegisterSignupServer(grpcServer, handler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}

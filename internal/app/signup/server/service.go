/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package server

import (
	"fmt"
	"github.com/nalej/grpc-application-go"
	"github.com/nalej/grpc-infrastructure-go"
	"net"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
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
	clusterClient grpc_infrastructure_go.ClustersClient
	appClient grpc_application_go.ApplicationsClient
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
	cClient := grpc_infrastructure_go.NewClustersClient(smConn)
	aClient := grpc_application_go.NewApplicationsClient(smConn)
	log.Debug().Str("smConn", smConn.GetState().String()).Str("uConn", uConn.GetState().String()).Msg("connections have been created")

	return &Clients{oClient, uClient, cClient, aClient}, nil
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

	manager := signup.NewManager(clients.orgClient, clients.userClient, clients.clusterClient, clients.appClient)
	handler := signup.NewHandler(manager, s.Configuration.UsePresharedSecret, s.Configuration.PresharedSecret)

	var grpcServer *grpc.Server
	if s.Configuration.UseTLS {
		creds, err := s.Configuration.GetTLSConfig()
		if err != nil {
			log.Fatal().Str("err", err.DebugReport()).Msg("error getting TLS configuration")
		}
		authData := AuthData{
			ClientSecret: s.Configuration.ClientSecret,
		}
		log.Debug().Msg("Creating server with TLS config")
		grpcServer = grpc.NewServer(
			grpc.Creds(creds),
			grpc.StreamInterceptor(grpc_auth.StreamServerInterceptor(authData.Authenticate)),
			grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(authData.Authenticate)),
		)
	} else {
		log.Debug().Msg("Creating server without certs")
		grpcServer = grpc.NewServer()
	}
	grpc_signup_go.RegisterSignupServer(grpcServer, handler)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	log.Info().Int("port", s.Configuration.Port).Msg("Launching gRPC server")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal().Errs("failed to serve: %v", []error{err})
	}
	return nil
}

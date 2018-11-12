/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package server

import (
	"github.com/nalej/derrors"
	"github.com/nalej/signup/version"
	"github.com/rs/zerolog/log"
)

//Config struct to hold the necessary service configuration
type Config struct {
	// Port where the gRPC API service will listen requests.
	Port int
	// HTTPPort where the HTTP gRPC gateway will be listening.
	HTTPPort int
	// SystemModelAddress with the host:port to connect to System Model
	SystemModelAddress string
	// UserManagerAddress with the host:port to connect to the User manager.
	UserManagerAddress string
}

//Validate makes the necessary validation in configuration prior to its use
func (conf *Config) Validate() derrors.Error {

	if conf.Port <= 0 {
		return derrors.NewInvalidArgumentError("ports must be valid")
	}

	if conf.SystemModelAddress == "" {
		return derrors.NewInvalidArgumentError("systemModelAddress must be set")
	}

	if conf.UserManagerAddress == "" {
		return derrors.NewInvalidArgumentError("userManagerAddress must be set")
	}

	return nil
}

//Print outputs the current configuration
func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC port")
	log.Info().Int("port", conf.HTTPPort).Msg("HTTP port")
	log.Info().Str("URL", conf.SystemModelAddress).Msg("System Model")
	log.Info().Str("URL", conf.UserManagerAddress).Msg("User Manager")
}

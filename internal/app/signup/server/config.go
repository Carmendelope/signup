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

package server

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"strings"

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
	// UseTLS if the gRPC service uses TLS or not
	UseTLS bool
	// CertCA with the absolute path to the certificate CA to trust
	CertCAPath string
	// CertFilePath with the absolute path to the certificate file for gRPC Server TLS
	CertFilePath string
	// CertKeyPath with the absolute path to the certificate key for gRPC Server TLS
	CertKeyPath string
	// ClientSecret with the client secret expected in client certificates
	ClientSecret string

	UsePresharedSecret bool
	PresharedSecret string
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
	if err := conf.validateTLS(); err != nil {
		return err
	}

	if conf.UsePresharedSecret && conf.PresharedSecret == "" {
		return derrors.NewInvalidArgumentError("preshared secret must be set")
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
	log.Info().Bool("TLS", conf.UseTLS).Msg("TLS Enabled")
	if conf.UseTLS{
		log.Info().Str("TLS", conf.CertCAPath).Msg("CA Certificate Path")
		log.Info().Str("TLS", conf.CertFilePath).Msg("Server Certificate Path")
		log.Info().Str("TLS", conf.CertKeyPath).Msg("Server Certificate Key Path")
		log.Info().Str("TLS", strings.Repeat("*", len(conf.ClientSecret))).Msg("Client certificate secret")
	}
	log.Info().Bool("enabled", conf.UsePresharedSecret).Msg("Use preshared secret")
	if conf.UsePresharedSecret{
		log.Info().Str("TLS", strings.Repeat("*", len(conf.PresharedSecret))).Msg("Preshared secret")
	}

}

func (conf *Config) validateTLS() derrors.Error {
	if conf.UseTLS {
		if conf.CertFilePath == "" || conf.CertKeyPath == "" {
			return derrors.NewInvalidArgumentError("if useTLS is enabled, certClientPath and certKeyPath must be set")
		}
		if _, err := tls.LoadX509KeyPair(conf.CertFilePath, conf.CertKeyPath); err != nil {
			return derrors.NewInvalidArgumentError("certFilePath or certKeyPath are invalid certificate file paths")
		}
	}
	return nil
}

//GetTLSConfig returns the necessary configuration with the CA and certificate files loaded if necessary
func (conf *Config) GetTLSConfig() (credentials.TransportCredentials, derrors.Error) {
	if conf.UseTLS {
		rootCAs := x509.NewCertPool()

		if conf.CertCAPath != "" {
			caCert, err := ioutil.ReadFile(conf.CertCAPath)
			if err != nil {
				log.Fatal().Errs("error loading CA certificate: %v", []error{err})
			}
			rootCAs.AppendCertsFromPEM(caCert)
		}

		serverCert, err := tls.LoadX509KeyPair(conf.CertFilePath, conf.CertKeyPath)
		if err != nil {
			log.Fatal().Errs("error loading Server certificate and key: %v", []error{err})
		}

		tlsConfig := &tls.Config{
			RootCAs:      rootCAs,
			ClientCAs:    rootCAs,
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{serverCert},
		}

		return credentials.NewTLS(tlsConfig), nil
	}
	return nil, derrors.NewGenericError("Requested TLS config without TLS enabled")
}

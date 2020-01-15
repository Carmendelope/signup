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

package commands

import (
	"io/ioutil"

	"github.com/nalej/signup/internal/app/signup/server"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config = server.Config{}

var clientSecretPath string

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launch the server API",
	Long:  `Launch the server API`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if config.UseTLS {
			contents, err := ioutil.ReadFile(clientSecretPath)
			if err != nil {
				panic(err)
			}
			config.ClientSecret = string(contents)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Launching API!")
		server := server.NewService(config)
		server.Run()
	},
}

func init() {
	runCmd.Flags().IntVar(&config.Port, "port", 8180, "Port to launch the Public gRPC API")
	runCmd.Flags().BoolVar(&config.UseTLS, "tls", false, "Enable TLS for gRPC Service")
	runCmd.Flags().StringVar(&config.CertCAPath, "caPath", "", "Absolute path to CA certificate")
	runCmd.Flags().StringVar(&config.CertFilePath, "certFilePath", "", "Absolute path to certificate file")
	runCmd.Flags().StringVar(&config.CertKeyPath, "certKeyPath", "", "Absolute path to certificate key")
	runCmd.Flags().StringVar(&clientSecretPath, "clientSecretPath", "", "Absolute path to client certificate secret")
	runCmd.PersistentFlags().StringVar(&config.SystemModelAddress, "systemModelAddress", "localhost:8800",
		"System Model address (host:port)")
	runCmd.PersistentFlags().StringVar(&config.UserManagerAddress, "userManagerAddress", "localhost:8920",
		"User Manager address (host:port)")
	runCmd.PersistentFlags().StringVar(&config.OrganizationManagerAddress, "organizationManagerAddress", "localhost:8950",
		"User Manager address (host:port)")
	runCmd.PersistentFlags().BoolVar(&config.UsePresharedSecret, "usePresharedSecret", false, "Use preshared secret to authenticate users")
	runCmd.PersistentFlags().StringVar(&config.PresharedSecret, "presharedSecret", "changemeifyouareusingthis", "Preshared secret with the client")
	rootCmd.AddCommand(runCmd)
}

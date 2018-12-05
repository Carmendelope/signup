/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
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
	runCmd.PersistentFlags().BoolVar(&config.UsePresharedSecret, "usePresharedSecret", false, "Use preshared secret to authenticate users")
	runCmd.PersistentFlags().StringVar(&config.PresharedSecret, "presharedSecret", "changemeifyouareusingthis", "Preshared secret with the client")
	rootCmd.AddCommand(runCmd)
}

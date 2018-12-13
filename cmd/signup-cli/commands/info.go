/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/signup/internal/app/cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Obtain the information of an existing organization",
	Long:  `Obtain the information of an existing organization`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		signupCli, err := cli.NewSignupCli(signupAddress, caPath, clientCertPath, clientKeyPath, presharedSecret)
		if err != nil {
			log.Error().Str("err", err.DebugReport()).Msg("cannot create CLI")
		}
		signupCli.Info(organizationID)
	},
}

func init() {
	infoCmd.Flags().StringVar(&organizationID, "organizationID", "", "Organization identifier")
	rootCmd.AddCommand(infoCmd)
}

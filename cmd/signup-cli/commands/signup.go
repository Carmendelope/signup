/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/signup/internal/app/cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var orgName string
var ownerEmail string
var ownerName string
var ownerPassword string

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Signup a new organization",
	Long:  `Signup a new organization creating the default roles and first user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		signupCli, err := cli.NewSignupCli(signupAddress)
		if err != nil{
			log.Error().Str("err", err.DebugReport()).Msg("cannot create CLI")
		}
		signupCli.SignupOrganization(orgName, ownerEmail, ownerName, ownerPassword)
	},
}

func init() {
	signupCmd.Flags().StringVar(&orgName, "orgName", "", "Name of the organization")
	signupCmd.Flags().StringVar(&ownerEmail, "ownerEmail", "", "Email of the organization owner")
	signupCmd.Flags().StringVar(&ownerName, "ownerName", "", "Name the owner")
	signupCmd.Flags().StringVar(&ownerPassword, "ownerPassword", "", "Password for the owner account")

	rootCmd.AddCommand(signupCmd)
}

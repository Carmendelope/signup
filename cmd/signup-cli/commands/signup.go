/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/signup/internal/app/cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var signupCmd = &cobra.Command{
	Use:   "signup",
	Short: "Signup a new organization",
	Long:  `Signup a new organization creating the default roles, the Nalej Admin, and first organization user`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		signupCli, err := cli.NewSignupCli(signupAddress, caPath, clientCertPath, clientKeyPath, presharedSecret)
		if err != nil {
			log.Fatal().Str("err", err.DebugReport()).Msg("cannot create CLI")
		}
		signupCli.SignupOrganization(orgName, ownerEmail, ownerName, ownerPassword, nalejAdminEmail, nalejAdminName, nalejAdminPassword)
	},
}

func init() {
	signupCmd.Flags().StringVar(&orgName, "orgName", "", "Name of the organization")
	_ = signupCmd.MarkFlagRequired("orgName")
	signupCmd.Flags().StringVar(&ownerEmail, "ownerEmail", "", "Email of the organization owner")
	_ = signupCmd.MarkFlagRequired("ownerEmail")
	signupCmd.Flags().StringVar(&ownerName, "ownerName", "", "Name the owner")
	_ = signupCmd.MarkFlagRequired("ownerName")
	signupCmd.Flags().StringVar(&ownerPassword, "ownerPassword", "", "Password for the owner account")
	_ = signupCmd.MarkFlagRequired("ownerPassword")
	signupCmd.Flags().StringVar(&nalejAdminEmail, "nalejAdminEmail", "", "Email of the Nalej administrator assigned to the organization")
	_ = signupCmd.MarkFlagRequired("nalejAdminEmail")
	signupCmd.Flags().StringVar(&nalejAdminName, "nalejAdminName", "", "Name the Nalej administrator")
	_ = signupCmd.MarkFlagRequired("nalejAdminName")
	signupCmd.Flags().StringVar(&nalejAdminPassword, "nalejAdminPassword", "", "Password for the Nalej administrator account")
	_ = signupCmd.MarkFlagRequired("nalejAdminPassword")
	rootCmd.AddCommand(signupCmd)
}

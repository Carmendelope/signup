/*
 * Copyright 2020 Nalej
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
			return
		}
		err = signupCli.SignupOrganization(orgName, orgEmail, orgFullAddress, orgCity, orgState, orgCountry, orgZipCode, orgPhotoPath,
			ownerEmail, ownerName, ownerLastName, ownerTitle, ownerPassword,
			nalejAdminEmail, nalejAdminName, nalejAdminLastName, nalejAdminTitle, nalejAdminPassword)
		if err != nil {
			log.Fatal().Msg("signup failed")
		}
	},
}

func init() {
	addOrgFlags()
	addOwnerFlags()
	addNalejAdminFlags()
}

// addOrgFlags adds the Organization fields as parameters for the command. Notice that default values are used
// to avoid breaking calling scripts.
func addOrgFlags() {
	signupCmd.Flags().StringVar(&orgName, "orgName", "", "Name of the organization")
	_ = signupCmd.MarkFlagRequired("orgName")
	// TODO orgEmail, orgAddress, orgCity, orgState, orgCountry, and orgZipCode must be marked as required when generation scripts are updated
	signupCmd.Flags().StringVar(&orgEmail, "orgEmail", "unknown@unknown.com", "Email of the organization")
	signupCmd.Flags().StringVar(&orgFullAddress, "orgAddress", "Unknown", "Organization full address")
	signupCmd.Flags().StringVar(&orgCity, "orgCity", "Unknown", "Organization City")
	signupCmd.Flags().StringVar(&orgState, "orgState", "Unknown", "Organization State")
	signupCmd.Flags().StringVar(&orgCountry, "orgCountry", "Unknown", "Organization State")
	signupCmd.Flags().StringVar(&orgZipCode, "orgZipCode", "Unknown", "Organization ZIP code")
	signupCmd.Flags().StringVar(&orgPhotoPath, "orgPhotoPath", "", "Path of the organization photo/logo")
}

// addOwnerFlags adds the Owner fields as parameters for the command. Notice that default values are used
//// to avoid breaking calling scripts.
func addOwnerFlags() {
	// TODO ownerLastName and ownerTitle must be marked as required when generation scripts are updated
	signupCmd.Flags().StringVar(&ownerEmail, "ownerEmail", "", "Email of the organization owner")
	_ = signupCmd.MarkFlagRequired("ownerEmail")
	signupCmd.Flags().StringVar(&ownerName, "ownerName", "", "Name the owner")
	_ = signupCmd.MarkFlagRequired("ownerName")
	signupCmd.Flags().StringVar(&ownerLastName, "ownerLastName", "Unknown", "Last name of the owner")
	signupCmd.Flags().StringVar(&ownerTitle, "ownerTitle", "Unknown", "Title of the owner")
	signupCmd.Flags().StringVar(&ownerPassword, "ownerPassword", "", "Password for the owner account")
	_ = signupCmd.MarkFlagRequired("ownerPassword")
}

// addNalejAdminFlags adds the Nalej Admin fields as parameters for the command. Notice that default values are used
//// to avoid breaking calling scripts.
func addNalejAdminFlags() {
	// TODO nalejAdminLastName and nalejAdminTitle must be marked as required when generation scripts are updated
	signupCmd.Flags().StringVar(&nalejAdminEmail, "nalejAdminEmail", "", "Email of the Nalej administrator assigned to the organization")
	_ = signupCmd.MarkFlagRequired("nalejAdminEmail")
	signupCmd.Flags().StringVar(&nalejAdminName, "nalejAdminName", "", "Name of the Nalej administrator")
	_ = signupCmd.MarkFlagRequired("nalejAdminName")
	signupCmd.Flags().StringVar(&nalejAdminLastName, "nalejAdminLastName", "Unknown", "Last name of the Nalej administrator")
	signupCmd.Flags().StringVar(&nalejAdminTitle, "nalejAdminTitle", "Unknown", "Title of the Nalej administrator")
	signupCmd.Flags().StringVar(&nalejAdminPassword, "nalejAdminPassword", "", "Password for the Nalej administrator account")
	_ = signupCmd.MarkFlagRequired("nalejAdminPassword")
	rootCmd.AddCommand(signupCmd)
}

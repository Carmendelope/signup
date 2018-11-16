/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/signup/version"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"fmt"
	"os"
)

var debugLevel bool
var consoleLogging bool

var signupAddress string
var caPath string
var clientCertPath string
var clientKeyPath string

var rootCmd = &cobra.Command{
	Use:     "signup-cli",
	Short:   "Signup CLI",
	Long:    `A command line tool to interact with the signup component`,
	Version: "unknown-version",
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugLevel, "debug", false, "Set debug level")
	rootCmd.PersistentFlags().BoolVar(&consoleLogging, "consoleLogging", false, "Pretty print logging")
	rootCmd.PersistentFlags().StringVar(&signupAddress, "signupAddress", "localhost:8180", "Signup address (host:port)")
	rootCmd.PersistentFlags().StringVar(&caPath, "caPath", "", "CA Certificate to use")
	rootCmd.PersistentFlags().StringVar(&clientCertPath, "clientCertPath", "", "Client certificate path")
	rootCmd.PersistentFlags().StringVar(&clientKeyPath, "clientKeyPath", "", "Client certificate key path")
}

// Execute runs the cli command execution chain
func Execute() {
	rootCmd.SetVersionTemplate(version.GetVersionInfo())
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// SetupLogging sets the debug level and console logging if required.
func SetupLogging() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugLevel {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if consoleLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

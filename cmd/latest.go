package cmd

import (
	"fmt"
	"os"
	"slices"

	"github.com/fun7257/sgv/internal/env"
	"github.com/fun7257/sgv/internal/installer"
	"github.com/fun7257/sgv/internal/version"
	"github.com/spf13/cobra"
)

var latestCmd = &cobra.Command{
	Use:   "latest",
	Short: "Install the latest Go version",
	Long:  `Check for the latest Go version, install it if not already installed, and switch to it.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the latest version
		latestVersion, err := version.GetLatestGoVersion()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting the latest Go version: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("The latest Go version is: %s\n", latestVersion)

		// Get the current version
		currentVersion := ""       // Initialize currentVersion
		currentVersionErr := false // Flag to indicate if there was an error getting current version

		cv, err := version.GetCurrentVersion()
		if err != nil {
			// No current version, this is expected for a fresh install
			fmt.Println("No Go version is currently active. Installing the latest version.")
			currentVersionErr = true // Set flag, but don't exit
			currentVersion = ""      // Set currentVersion to empty string
		} else {
			currentVersion = cv
		}

		if !currentVersionErr && currentVersion == latestVersion {
			fmt.Printf("You are already using the latest Go version: %s\n", latestVersion)
			return
		}

		// Check if the latest version is already installed
		localVersions, err := version.GetLocalVersions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting local Go versions: %v\n", err)
			os.Exit(1)
		}

		isInstalled := slices.Contains(localVersions, latestVersion)

		if isInstalled {
			fmt.Printf("Go version %s is already installed.\n", latestVersion)
		} else {
			fmt.Printf("Go version %s not found locally. Installing...\n", latestVersion)
			if err := installer.Install(latestVersion); err != nil {
				fmt.Fprintf(os.Stderr, "Error installing Go version %s: %v\n", latestVersion, err)
				os.Exit(1)
			}
		}

		// Switch to the latest version
		if err := version.SwitchToVersion(latestVersion); err != nil {
			fmt.Fprintf(os.Stderr, "Error switching to Go version %s: %v\n", latestVersion, err)
			os.Exit(1)
		}
		fmt.Printf("Successfully switched to Go version %s\n", latestVersion)

		// Check for and notify about environment variables
		if envVars, err := env.LoadEnvVars(latestVersion); err == nil && len(envVars) > 0 {
			fmt.Printf("Loading %d custom environment variables for %s...\n", len(envVars), latestVersion)
		}
	},
}

func init() {
	rootCmd.AddCommand(latestCmd)
}

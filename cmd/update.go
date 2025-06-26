package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fun7257/sgv/internal/installer"
	"github.com/fun7257/sgv/internal/version"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
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
		currentVersion, err := version.GetCurrentVersion()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting the current Go version: %v\n", err)
			os.Exit(1)
		}

		if currentVersion == latestVersion {
			fmt.Printf("You are already using the latest Go version: %s\n", latestVersion)
			return
		}

		// The version from the website has a 'go' prefix, remove it for installation.
		installVersion := strings.TrimPrefix(latestVersion, "go")

		// Check if the latest version is already installed
		localVersions, err := version.GetLocalVersions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting local Go versions: %v\n", err)
			os.Exit(1)
		}

		isInstalled := false
		for _, v := range localVersions {
			if v == latestVersion {
				isInstalled = true
				break
			}
		}

		if isInstalled {
			fmt.Printf("Go version %s is already installed.\n", latestVersion)
		} else {
			fmt.Printf("Go version %s not found locally. Installing...\n", latestVersion)
			if err := installer.Install(installVersion); err != nil {
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
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

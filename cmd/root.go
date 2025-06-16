package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fun7257/sgv/internal/config"
	"github.com/fun7257/sgv/internal/installer"
	"github.com/fun7257/sgv/internal/version"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sgv",
	Short: "SGV is a Go Version manager",
	Long: `A fast and flexible Go Version manager built with love by Howell.

This tool allows you to easily install and switch between different Go versions.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		versionStr := args[0]

		// Normalize version string (e.g., "1.22.1" -> "go1.22.1")
		if !strings.HasPrefix(versionStr, "go") {
			versionStr = "go" + versionStr
		}

		// Check if version is already installed
		installPath := filepath.Join(config.VERSIONS_DIR, versionStr)
		if _, err := os.Stat(installPath); os.IsNotExist(err) {
			fmt.Printf("Go version %s not found locally. Installing...\n", versionStr)
			if err := installer.Install(versionStr); err != nil {
				fmt.Fprintf(os.Stderr, "Error installing Go version %s: %v\n", versionStr, err)
				os.Exit(1)
			}
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking installation path: %v\n", err)
			os.Exit(1)
		}

		// Switch to the specified version
		if err := version.SwitchToVersion(versionStr); err != nil {
			fmt.Fprintf(os.Stderr, "Error switching to Go version %s: %v\n", versionStr, err)
			os.Exit(1)
		}

		fmt.Printf("Successfully switched to Go version %s\n", versionStr)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(config.Init)
}

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
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

		// Check if the requested version is supported
		if !isGoVersionSupported(versionStr) {
			fmt.Fprintf(os.Stderr, "Error: Go version %s is not supported. sgv only supports Go 1.13 and later.\n", versionStr)
			os.Exit(1)
		}

		goModVersion, err := findGoModVersion()
		if err != nil {
			// If there's an error finding go.mod, it means it's not a Go project or an error occurred.
			// We will not exit, but continue with the user's requested version.
			fmt.Fprintf(os.Stderr, "Warning: Could not determine go.mod version: %v\n", err)
			goModVersion = "" // Ensure goModVersion is empty to skip go.mod related logic
		}

		if goModVersion != "" {
			// Compare versions using semver
			if !isGoVersionCompatible(versionStr, goModVersion) {
				fmt.Fprintf(os.Stderr, "Error: The requested Go version %s is lower than the go.mod requirement %s. Please switch to a compatible version manually.\n", versionStr, goModVersion)
				os.Exit(1)
			}
		}

		// Check if version is already installed
		installPath := filepath.Join(config.VersionsDir, versionStr)
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
	checkPlatformSupport()
	cobra.OnInitialize(config.Init)
}

// checkPlatformSupport ensures the current platform is supported
func checkPlatformSupport() {
	if runtime.GOOS == "windows" {
		fmt.Fprintf(os.Stderr, "Error: Windows is not supported by sgv. This tool only works on macOS and Linux.\n")
		os.Exit(1)
	}
}

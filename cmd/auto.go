package cmd

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/fun7257/sgv/internal/version"

	"github.com/spf13/cobra"
)

var autoCmd = &cobra.Command{
	Use:   "auto",
	Short: "Automatically switch to the most suitable Go version for the current project",
	Long:  `If the current directory is a Go project, this command automatically switches to the Go version specified in go.mod, or the closest compatible installed version. If no compatible version is found, it prompts the user to install one.`,
	Run: func(cmd *cobra.Command, args []string) {
		goModVersion, err := findGoModVersion()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error determining go.mod version: %v\n", err)
			os.Exit(1)
		}

		if goModVersion == "" {
			fmt.Println("Current directory is not a Go project (no go.mod found).")
			return
		}

		// Check if the go.mod version is supported
		if !isGoVersionSupported(goModVersion) {
			fmt.Fprintf(os.Stderr, "Error: The Go version required by go.mod (%s) is not supported. sgv only supports Go 1.13 and later.\n", goModVersion)
			os.Exit(1)
		}

		localVersions, err := version.GetLocalVersions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting local versions: %v\n", err)
			os.Exit(1)
		}

		suitableVersion := goModVersion
		suitableVersionSource := "remote" // Assume remote initially

		// Check if the go.mod version is installed locally
		if slices.Contains(localVersions, goModVersion) {
			suitableVersionSource = "local"
		}

		currentActiveVersion, err := version.GetCurrentVersion()
		if err != nil {
			// If we can't get current version, proceed with suitableVersion
			currentActiveVersion = ""
		}

		if suitableVersion != "" {
			// If suitableVersion is the same as currentActiveVersion, no switch needed.
			if currentActiveVersion != "" && suitableVersion == currentActiveVersion {
				return // No output, no switch needed
			}

			fmt.Printf("go.mod requires Go version: %s\n", goModVersion)
			msg := fmt.Sprintf("Found suitable version: %s.", suitableVersion)
			if suitableVersionSource == "remote" {
				msg += " (Will download and install)"
			}
			fmt.Println(msg)
			fmt.Printf("Switch to this version? (y/n): ")

			var response string
			_, err := fmt.Scanln(&response)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Invalid input: %v\n", err)
				os.Exit(1)
			}

			if strings.ToLower(response) == "y" {
				rootCmd.SetArgs([]string{suitableVersion})
				if err := rootCmd.Execute(); err != nil {
					fmt.Fprintf(os.Stderr, "Error switching to Go version %s: %v\n", suitableVersion, err)
					os.Exit(1)
				}
			} else {
				fmt.Println("Switch aborted.")
			}
		} else {
			fmt.Printf("No Go version found (local or remote) that meets the go.mod requirement (%s). Please install a compatible version manually.\n", goModVersion)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(autoCmd)
}

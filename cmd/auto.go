package cmd

import (
	"fmt"
	"os"

	"github.com/fun7257/sgv/internal/version"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
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

		fmt.Printf("go.mod requires Go version: %s\n", goModVersion)

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

		var suitableVersion string
		// Find the smallest installed version that is >= goModVersion
		for _, lv := range localVersions {
			if isGoVersionSupported(lv) && isGoVersionCompatible(lv, goModVersion) {
				if suitableVersion == "" || semver.Compare(normalizeGoVersion(lv), normalizeGoVersion(suitableVersion)) < 0 {
					suitableVersion = lv
				}
			}
		}

		if suitableVersion != "" {
			currentActiveVersion, err := version.GetCurrentVersion()
			if err == nil && currentActiveVersion == suitableVersion {
				fmt.Printf("Go version %s is already the most suitable version for this project.\n", suitableVersion)
				return
			}

			currentActiveVersion, err = version.GetCurrentVersion()
			if err == nil && currentActiveVersion != "" {
				// If current active version is newer than the suitable version, don't switch
				if semver.Compare(normalizeGoVersion(currentActiveVersion), normalizeGoVersion(suitableVersion)) > 0 {
					fmt.Printf("Current Go version %s is newer than the go.mod requirement. No switch needed.\n", currentActiveVersion)
					return
				}
			}

			fmt.Printf("Found suitable installed version: %s. Switching...\n", suitableVersion)
			if err = version.SwitchToVersion(suitableVersion); err != nil {
				fmt.Fprintf(os.Stderr, "Error switching to Go version %s: %v\n", suitableVersion, err)
				os.Exit(1)
			}
			fmt.Printf("Successfully switched to Go version %s.\n", suitableVersion)
		} else {
			fmt.Printf("No installed Go version found that meets the go.mod requirement (%s). Please install a compatible version.\n", goModVersion)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(autoCmd)
}

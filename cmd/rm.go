package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fun7257/sgv/internal/config"
	"github.com/fun7257/sgv/internal/version"
	"github.com/samber/lo"

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm [version...]",
	Short: "Uninstall one or more Go versions",
	Long: `Uninstall one or more previously installed Go versions from your system.

You can specify multiple full version numbers (e.g., 1.22.1) or major versions (e.g., 1.22) to remove all its sub-versions.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		currentVersion, err := version.GetCurrentVersion()
		if err != nil {
			// It's not critical if we can't get the current version, but we should inform the user.
			fmt.Fprintf(os.Stderr, "Warning: could not determine current Go version: %v\n", err)
		}

		installedVersions, err := version.GetLocalVersions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: could not list installed Go versions: %v\n", err)
			os.Exit(1)
		}
		installedSet := lo.SliceToMap(installedVersions, func(v string) (string, struct{}) {
			return v, struct{}{}
		})

		versionsToUninstall := findVersionsToUninstall(args, installedVersions)

		// Filter out non-existent versions and the active version
		var finalVersionsToUninstall []string
		for _, v := range versionsToUninstall {
			if _, exists := installedSet[v]; !exists {
				fmt.Fprintf(os.Stderr, "Info: Go version %s is not installed. Skipping.\n", v)
				continue
			}
			if v == currentVersion {
				fmt.Fprintf(os.Stderr, "Info: Cannot uninstall currently active Go version (%s). It will be skipped.\n", v)
				continue
			}
			finalVersionsToUninstall = append(finalVersionsToUninstall, v)
		}

		// Remove duplicates
		finalVersionsToUninstall = lo.Uniq(finalVersionsToUninstall)

		if len(finalVersionsToUninstall) == 0 {
			fmt.Println("No versions to uninstall.")
			return
		}

		// Confirmation prompt
		fmt.Printf("Are you sure you want to uninstall the following Go versions: %s? (y/N): ", strings.Join(finalVersionsToUninstall, ", "))
		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			fmt.Println("Uninstallation cancelled.")
			return
		}

		// Delete the version directories
		for _, v := range finalVersionsToUninstall {
			fmt.Printf("Uninstalling Go version %s...\n", v)
			versionPath := filepath.Join(config.VersionsDir, v)
			if err := os.RemoveAll(versionPath); err != nil {
				fmt.Fprintf(os.Stderr, "Error uninstalling Go version %s: %v\n", v, err)
				// Continue to the next version instead of exiting
			} else {
				fmt.Printf("Successfully uninstalled Go version %s.\n", v)
			}
		}
	},
}

// findVersionsToUninstall resolves user input (like "1.22") into a list of full version strings.
func findVersionsToUninstall(args []string, installedVersions []string) []string {
	var versionsToUninstall []string
	for _, arg := range args {
		normalizedArg := arg
		if !strings.HasPrefix(normalizedArg, "go") {
			normalizedArg = "go" + normalizedArg
		}

		// Check if it's a major version (e.g., "go1.22")
		isMajorVersion := !strings.Contains(strings.TrimPrefix(normalizedArg, "go"), ".")
		if !isMajorVersion {
			parts := strings.Split(strings.TrimPrefix(normalizedArg, "go"), ".")
			if len(parts) == 2 {
				isMajorVersion = true
			}
		}

		if isMajorVersion {
			// It's a major version, find all installed sub-versions
			found := false
			for _, installed := range installedVersions {
				if strings.HasPrefix(installed, normalizedArg+".") || installed == normalizedArg {
					versionsToUninstall = append(versionsToUninstall, installed)
					found = true
				}
			}
			if !found {
				fmt.Fprintf(os.Stderr, "Info: No installed versions found for major version %s.\n", arg)
			}
		} else {
			// It's a full version string
			versionsToUninstall = append(versionsToUninstall, normalizedArg)
		}
	}
	return versionsToUninstall
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

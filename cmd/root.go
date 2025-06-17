package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fun7257/sgv/internal/config"
	"github.com/fun7257/sgv/internal/installer"
	"github.com/fun7257/sgv/internal/version"

	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

// findGoModVersion searches for a go.mod file in the current directory or parent directories
// and returns the Go version specified in it.
func findGoModVersion() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	goModPath := filepath.Join(currentDir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		content, err := os.ReadFile(goModPath)
		if err != nil {
			return "", fmt.Errorf("failed to read go.mod file: %w", err)
		}

		re := regexp.MustCompile(`^go\s+(\d+\.\d+(\.\d+)?)$`)
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				return "go" + matches[1], nil
			}
		}
		return "", fmt.Errorf("go.mod found, but no go version specified")
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("error checking for go.mod file: %w", err)
	}

	return "", nil // No go.mod found in current directory
}

// normalizeGoVersion ensures the version string starts with 'v' for semver comparison.
func normalizeGoVersion(v string) string {
	if !strings.HasPrefix(v, "go") {
		v = "go" + v
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + strings.TrimPrefix(v, "go")
	}
	return v
}

// isGoVersionCompatible checks if candidateVersion is greater than or equal to requiredVersion.
func isGoVersionCompatible(candidateVersion, requiredVersion string) bool {
	candidate := normalizeGoVersion(candidateVersion)
	required := normalizeGoVersion(requiredVersion)
	return semver.Compare(candidate, required) >= 0
}

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
				fmt.Printf("Warning: The requested Go version %s is lower than the go.mod requirement %s.\n", versionStr, goModVersion)
				fmt.Println("Attempting to find a suitable installed version...")

				localVersions, err := version.GetLocalVersions()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting local versions: %v\n", err)
					os.Exit(1)
				}

				var suitableVersion string
				// Find the smallest installed version that is >= goModVersion
				for _, lv := range localVersions {
					if isGoVersionCompatible(lv, goModVersion) { // Use semver comparison
						if suitableVersion == "" || semver.Compare(normalizeGoVersion(lv), normalizeGoVersion(suitableVersion)) < 0 {
							suitableVersion = lv
						}
					}
				}

				if suitableVersion != "" {
					fmt.Printf("Switching to installed Go version %s, which meets the go.mod requirement.\n", suitableVersion)
					versionStr = suitableVersion // Override the user's requested version
				} else {
					fmt.Printf("No installed Go version found that meets the go.mod requirement (%s). Please install a compatible version.\n", goModVersion)
					os.Exit(1)
				}
			}
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

package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fun7257/sgv/internal/config"
	"github.com/fun7257/sgv/internal/version"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [version]",
	Short: "Uninstall a specific Go version",
	Long:  `Uninstall a previously installed Go version from your system.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		versionToUninstall := args[0]

		// Normalize version string (e.g., "1.22.1" -> "go1.22.1")
		if !strings.HasPrefix(versionToUninstall, "go") {
			versionToUninstall = "go" + versionToUninstall
		}

		// Check if the version to be uninstalled is the current active version
		currentVersion, err := version.GetCurrentVersion()
		if err == nil && currentVersion == versionToUninstall {
			fmt.Fprintf(os.Stderr, "Error: Cannot uninstall currently active Go version. Please switch to another version first.\n")
			os.Exit(1)
		}

		versionPath := filepath.Join(config.VersionsDir, versionToUninstall)

		// Check if the version directory actually exists
		if _, err := os.Stat(versionPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Go version %s is not installed.\n", versionToUninstall)
			os.Exit(1)
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking version path: %v\n", err)
			os.Exit(1)
		}

		// Confirmation prompt
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("Are you sure you want to uninstall Go version %s? (y/N): ", versionToUninstall)
		response, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(response)) != "y" {
			fmt.Println("Uninstallation cancelled.")
			return
		}

		// Delete the version's directory
		fmt.Printf("Uninstalling Go version %s...\n", versionToUninstall)
		if err := os.RemoveAll(versionPath); err != nil {
			fmt.Fprintf(os.Stderr, "Error uninstalling Go version %s: %v\n", versionToUninstall, err)
			os.Exit(1)
		}

		fmt.Printf("Successfully uninstalled Go version %s.\n", versionToUninstall)
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

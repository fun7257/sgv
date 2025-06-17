package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/fun7257/sgv/internal/installer"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [version]",
	Short: "Download and install a specific Go version",
	Long:  `Download and install a specific Go version to your system without switching to it.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		versionToInstall := args[0]

		// Normalize version string (e.g., "1.22.1" -> "go1.22.1")
		if !strings.HasPrefix(versionToInstall, "go") {
			versionToInstall = "go" + versionToInstall
		}

		fmt.Printf("Installing Go version %s...\n", versionToInstall)
		if err := installer.Install(versionToInstall); err != nil {
			fmt.Fprintf(os.Stderr, "Error installing Go version %s: %v\n", versionToInstall, err)
			os.Exit(1)
		}

		fmt.Printf("Successfully installed Go version %s.\n", versionToInstall)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}

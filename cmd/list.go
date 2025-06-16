package cmd

import (
	"fmt"
	"os"

	"github.com/fun7257/sgv/internal/version"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed Go versions",
	Long:  `List all Go versions installed by sgv, and indicate the currently active version.`,
	Run: func(cmd *cobra.Command, args []string) {
		localVersions, err := version.GetLocalVersions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting local versions: %v\n", err)
			os.Exit(1)
		}

		currentVersion, err := version.GetCurrentVersion()
		if err != nil {
			// If there's an error getting the current version, it might mean no version is active.
			// We'll just proceed without highlighting a current version.
			currentVersion = ""
		}

		if len(localVersions) == 0 {
			fmt.Println("No Go versions installed yet.")
			return
		}

		fmt.Println("Installed Go versions:")
		for _, v := range localVersions {
			if v == currentVersion {
				fmt.Printf("  %s %s\n", v, color.GreenString("<- current"))
			} else {
				fmt.Printf("  %s\n", v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

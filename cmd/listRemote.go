package cmd

import (
	"fmt"
	"os"

	"github.com/fun7257/sgv/internal/version"

	"github.com/spf13/cobra"
)

var listRemoteCmd = &cobra.Command{
	Use:   "list-remote",
	Short: "List available Go versions from remote",
	Long:  `List all Go versions available for download from the official Go website.`,
	Run: func(cmd *cobra.Command, args []string) {
		remoteVersions, err := version.GetRemoteVersions()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting remote versions: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("Available Go versions (remote):")
		for _, v := range remoteVersions {
			if v.Stable {
				fmt.Printf("  %s\n", v.Version)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listRemoteCmd)
}

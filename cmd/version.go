package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of SGV",
	Long:  `All software has versions. This is SGV's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Go Version: %s\n", goVersion)
		fmt.Printf("Commit: %s\n", commit)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

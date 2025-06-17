package cmd

import (
	"fmt"

	"github.com/fun7257/sgv/internal/version" // Import the version package
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of SGV",
	Long:  `All software has versions. This is SGV's.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SGV Version: %s\n", version.GetSGVVersion()) // Display SGV's version
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

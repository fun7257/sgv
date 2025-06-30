package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

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
		// Group versions by major version
		groupedVersions := make(map[string][]string)
		for _, v := range localVersions {
			majorVersion := getMajorVersion(v)
			groupedVersions[majorVersion] = append(groupedVersions[majorVersion], v)
		}

		// Get sorted keys
		var sortedKeys []string
		for k := range groupedVersions {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			fmt.Printf("%s:\n", k)
			for _, v := range groupedVersions[k] {
				if v == currentVersion {
					fmt.Printf("  %s %s\n", v, color.GreenString("<- current"))
				} else {
					fmt.Printf("  %s\n", v)
				}
			}
		}
	},
}

func getMajorVersion(version string) string {
	parts := strings.Split(strings.TrimPrefix(version, "go"), ".")
	if len(parts) >= 2 {
		return "go" + parts[0] + "." + parts[1]
	}
	return version
}

func init() {
	rootCmd.AddCommand(listCmd)
}

package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fun7257/sgv/internal/version"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

// subCmd represents the sub command
var subCmd = &cobra.Command{
	Use:   "sub [major_version]",
	Short: "List minor versions for a specific Go major version",
	Long:  `List all available minor patch versions for a given Go major version.
For example: sgv sub 1.22`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		majorVersion := strings.TrimPrefix(args[0], "go")
		if !strings.HasPrefix(majorVersion, "1.") {
			majorVersion = "1." + majorVersion
		}

		// Check if the major version is at least 1.13
		if semver.Compare("v"+majorVersion, "v1.13") < 0 {
			return fmt.Errorf("this command is only available for Go versions 1.13 and higher")
		}

		allVersions, err := version.FetchAllGoVersions()
		if err != nil {
			return fmt.Errorf("failed to fetch Go versions: %w", err)
		}

		localVersions, err := version.GetLocalVersions()
		if err != nil {
			return fmt.Errorf("failed to get local Go versions: %w", err)
		}

		localVersionSet := make(map[string]struct{})
		for _, v := range localVersions {
			localVersionSet[v] = struct{}{}
		}

		var matchedVersions []string
		for _, v := range allVersions {
			if strings.HasPrefix(v, "go"+majorVersion) {
				matchedVersions = append(matchedVersions, v)
			}
		}

		sort.Slice(matchedVersions, func(i, j int) bool {
			return semver.Compare("v"+strings.TrimPrefix(matchedVersions[i], "go"), "v"+strings.TrimPrefix(matchedVersions[j], "go")) < 0
		})

		fmt.Printf("Available minor versions for go%s:\n", majorVersion)
		if len(matchedVersions) == 0 {
			fmt.Println("No versions found for the specified major version.")
			return nil
		}

		for _, v := range matchedVersions {
			if _, ok := localVersionSet[v]; ok {
				fmt.Printf("%s (installed)\n", v)
			} else {
				fmt.Println(v)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(subCmd)
}
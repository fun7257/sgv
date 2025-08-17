package cmd

import (
	"fmt"
	"runtime"
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
	Long: `List all available minor patch versions for a given Go major version.
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

		// Get current platform info
		currentOS := runtime.GOOS
		currentArch := runtime.GOARCH

		allVersions, err := version.GetStableGoVersions()
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

		var matchedVersions []version.GoVersion
		for _, v := range allVersions {
			if strings.HasPrefix(v.Version, "go"+majorVersion) {
				matchedVersions = append(matchedVersions, v)
			}
		}

		// Group versions by version string and check platform compatibility
		versionMap := make(map[string][]version.GoVersion)
		for _, v := range matchedVersions {
			versionMap[v.Version] = append(versionMap[v.Version], v)
		}

		// Convert back to sorted slice with platform info
		var sortedVersions []string
		for versionStr := range versionMap {
			sortedVersions = append(sortedVersions, versionStr)
		}
		sort.Slice(sortedVersions, func(i, j int) bool {
			return semver.Compare("v"+strings.TrimPrefix(sortedVersions[i], "go"), "v"+strings.TrimPrefix(sortedVersions[j], "go")) < 0
		})

		fmt.Printf("Available minor versions for go%s:\n", majorVersion)
		if len(sortedVersions) == 0 {
			fmt.Println("No versions found for the specified major version.")
			return nil
		}

		// ANSI color codes
		const (
			colorReset = "\033[0m"
			colorGray  = "\033[90m"
		)

		for _, versionStr := range sortedVersions {
			platforms := versionMap[versionStr]

			// Check if current platform is supported
			isCurrentPlatformSupported := false
			for _, p := range platforms {
				if p.OS == currentOS && p.Arch == currentArch {
					isCurrentPlatformSupported = true
					break
				}
			}

			// Format output based on installation status and platform support
			if _, ok := localVersionSet[versionStr]; ok {
				// Already installed
				if isCurrentPlatformSupported {
					fmt.Printf("%s (installed)\n", versionStr)
				} else {
					fmt.Printf("%s%s (installed, incompatible)%s\n", colorGray, versionStr, colorReset)
				}
			} else {
				// Not installed
				if isCurrentPlatformSupported {
					fmt.Printf("%s\n", versionStr)
				} else {
					fmt.Printf("%s%s (incompatible with %s/%s)%s\n", colorGray, versionStr, currentOS, currentArch, colorReset)
				}
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(subCmd)
}

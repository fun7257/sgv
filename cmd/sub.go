package cmd

import (
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/fun7257/sgv/internal/version"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
)

var interactive bool

// subCmd represents the sub command
var subCmd = &cobra.Command{
	Use:   "sub [major_version]",
	Short: "List minor versions for a specific Go major version",
	Long: `List all available minor patch versions for a given Go major version.
  Example: sgv sub 1.22
  Use -i or --interactive flag to interactively select and install a version using arrow keys.`,
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

		// Store installable versions (for interactive mode)
		var installableVersions []string

		if interactive {
			// In interactive mode, just show a simple list and collect installable versions
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

				// Simple display for interactive mode
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
						installableVersions = append(installableVersions, versionStr)
					} else {
						fmt.Printf("%s%s (incompatible with %s/%s)%s\n", colorGray, versionStr, currentOS, currentArch, colorReset)
					}
				}
			}

			// Interactive selection mode
			return handleInteractiveSelection(sortedVersions, installableVersions, localVersionSet)
		} else {
			// Non-interactive mode, show detailed list
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
		}

		return nil
	},
}

// handleInteractiveSelection handles the interactive version selection with arrow keys
func handleInteractiveSelection(allVersions []string, installableVersions []string, localVersionSet map[string]struct{}) error {
	if len(installableVersions) == 0 {
		fmt.Println("\nNo installable versions available for your platform.")
		return nil
	}

	// Create selection items with status information - only include selectable items
	type VersionItem struct {
		Version     string
		DisplayName string
	}

	var selectableItems []VersionItem
	for _, versionStr := range allVersions {
		_, isInstalled := localVersionSet[versionStr]
		isInstallable := false
		for _, installable := range installableVersions {
			if installable == versionStr {
				isInstallable = true
				break
			}
		}

		// Only add installable (not installed and compatible) versions to the selector
		if !isInstalled && isInstallable {
			selectableItems = append(selectableItems, VersionItem{
				Version:     versionStr,
				DisplayName: versionStr,
			})
		}
	}

	if len(selectableItems) == 0 {
		fmt.Println("\nAll compatible versions are already installed.")
		return nil
	}

	// Create promptui selector
	prompt := promptui.Select{
		Label: "Select a Go version to install",
		Items: selectableItems,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}:",
			Active:   "▸ {{ .DisplayName | cyan }}",
			Inactive: "  {{ .DisplayName }}",
			Selected: "{{ .DisplayName | green }}",
		},
		Size: 10, // Show up to 10 items at once
	}

	fmt.Println() // Add some spacing
	idx, _, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrInterrupt {
			fmt.Println("Selection cancelled.")
			return nil
		}
		return fmt.Errorf("selection failed: %w", err)
	}

	selectedItem := selectableItems[idx]
	selectedVersion := selectedItem.Version

	// Confirm installation
	confirmPrompt := promptui.Prompt{
		Label:     fmt.Sprintf("Install and switch to %s", selectedVersion),
		IsConfirm: true,
	}

	_, err = confirmPrompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			fmt.Println("Installation cancelled.")
			return nil
		}
		return fmt.Errorf("confirmation failed: %w", err)
	}

	// Execute the root command to install and switch
	fmt.Printf("Installing %s...\n", selectedVersion)
	rootCmd.SetArgs([]string{selectedVersion})
	return rootCmd.Execute()
}

func init() {
	subCmd.Flags().BoolVarP(&interactive, "interactive", "i", false, "Interactive mode to select and install a version using arrow keys")
	rootCmd.AddCommand(subCmd)
}

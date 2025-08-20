package cmd

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fun7257/sgv/internal/env"
	"github.com/spf13/cobra"
)

var (
	writeFlag string
	unsetFlag string
	shellFlag bool
	cleanFlag bool
	clearFlag bool
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables for Go versions",
	Long: `Manage environment variables for specific Go versions.

Examples:
  sgv env                      # List all environment variables for current version
  sgv env -w GOWORK=auto      # Set GOWORK environment variable
  sgv env -u GODEBUG          # Remove GODEBUG environment variable
  sgv env --clear             # Clear all environment variables for current version
  sgv env --shell             # Output environment variables in shell format
  sgv env --shell --clean     # Output shell format with cleanup of old variables`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get current version
		currentVersion, err := env.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("cannot determine current Go version: %w", err)
		}

		// Handle clear operation
		if clearFlag {
			return handleClearFlag(currentVersion)
		}

		// Handle shell output format
		if shellFlag {
			return outputShellFormat(currentVersion)
		}

		// Handle write operation
		if writeFlag != "" {
			return handleWriteFlag(currentVersion, writeFlag)
		}

		// Handle unset operation
		if unsetFlag != "" {
			return handleUnsetFlag(currentVersion, unsetFlag)
		}

		// Default: list all environment variables
		return listEnvVars(currentVersion)
	},
}

func handleWriteFlag(version, writeValue string) error {
	// Parse KEY=VALUE
	parts := strings.SplitN(writeValue, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid format: expected KEY=VALUE, got %s", writeValue)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return fmt.Errorf("environment variable name cannot be empty")
	}

	// Set the environment variable
	if err := env.SetEnvVar(version, key, value); err != nil {
		return err
	}

	fmt.Printf("Environment variable %s set to '%s' for Go version %s\n", key, value, version)
	return nil
}

func handleUnsetFlag(version, key string) error {
	key = strings.TrimSpace(key)
	if key == "" {
		return fmt.Errorf("environment variable name cannot be empty")
	}

	// Remove the environment variable
	if err := env.UnsetEnvVar(version, key); err != nil {
		return err
	}

	fmt.Printf("Environment variable %s removed for Go version %s\n", key, version)
	return nil
}

func handleClearFlag(version string) error {
	// Load existing variables to show what will be cleared
	vars, err := env.LoadEnvVars(version)
	if err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	if len(vars) == 0 {
		fmt.Printf("No environment variables set for Go version %s\n", version)
		return nil
	}

	// Show what will be cleared
	fmt.Printf("The following environment variables will be cleared for Go version %s:\n", version)
	keys := make([]string, 0, len(vars))
	for key := range vars {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		fmt.Printf("  %s=%s\n", key, vars[key])
	}

	// Confirm with user
	fmt.Print("Are you sure you want to clear all these variables? (y/N): ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(strings.TrimSpace(response)) != "y" {
		fmt.Println("Operation cancelled.")
		return nil
	}

	// Clear all variables
	if err := env.ClearAllEnvVars(version); err != nil {
		return err
	}

	fmt.Printf("All environment variables cleared for Go version %s\n", version)
	return nil
}

func listEnvVars(version string) error {
	// Load environment variables
	vars, err := env.LoadEnvVars(version)
	if err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	fmt.Printf("Current Go version: %s\n", version)

	if len(vars) == 0 {
		fmt.Println("No custom environment variables set.")
		return nil
	}

	fmt.Println("Environment variables:")

	// Sort keys for consistent output
	keys := make([]string, 0, len(vars))
	for key := range vars {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Print variables
	for _, key := range keys {
		fmt.Printf("  %s=%s\n", key, vars[key])
	}

	return nil
}

func outputShellFormat(version string) error {
	// Load environment variables for the current version first
	currentVars, err := env.LoadEnvVars(version)
	if err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Clean environment variables if --clean flag is specified
	if cleanFlag {
		// Get all environment variables from all versions
		allVars, err := env.GetAllEnvVars()
		if err != nil {
			return fmt.Errorf("failed to get all environment variables: %w", err)
		}

		// Collect variables from other versions (not current version)
		otherVersionVars := make(map[string]bool)
		for ver, vars := range allVars {
			if ver != version { // Skip current version
				for key := range vars {
					otherVersionVars[key] = true
				}
			}
		}

		// Unset variables that exist in other versions but not in current version
		for key := range otherVersionVars {
			// Skip protected variables
			if !env.IsProtectedVar(key) {
				// Only unset if this variable is NOT in the current version
				if _, existsInCurrent := currentVars[key]; !existsInCurrent {
					fmt.Printf("unset %s\n", key)
				}
			}
		}
	}

	// Set the current version's environment variables
	for key, value := range currentVars {
		fmt.Printf("export %s='%s'\n", key, value)
	}

	return nil
}

func init() {
	// Add flags
	envCmd.Flags().StringVarP(&writeFlag, "write", "w", "", "Set environment variable (format: KEY=VALUE)")
	envCmd.Flags().StringVarP(&unsetFlag, "unset", "u", "", "Remove environment variable")
	envCmd.Flags().BoolVar(&shellFlag, "shell", false, "Output environment variables in shell format")
	envCmd.Flags().BoolVar(&cleanFlag, "clean", false, "Clean (unset) all environment variables before setting new ones (only works with --shell)")
	envCmd.Flags().BoolVar(&clearFlag, "clear", false, "Clear all environment variables for current version")

	// Make flags mutually exclusive (except clean can be used with shell)
	envCmd.MarkFlagsMutuallyExclusive("write", "unset", "clear")

	rootCmd.AddCommand(envCmd)
}

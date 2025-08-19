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
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage environment variables for Go versions",
	Long: `Manage environment variables for specific Go versions.

Examples:
  sgv env                      # List all environment variables for current version
  sgv env -w GOWORK=auto      # Set GOWORK environment variable
  sgv env -u GODEBUG          # Remove GODEBUG environment variable
  sgv env --shell             # Output environment variables in shell format`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get current version
		currentVersion, err := env.GetCurrentVersion()
		if err != nil {
			return fmt.Errorf("cannot determine current Go version: %w", err)
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
	// Load environment variables
	vars, err := env.LoadEnvVars(version)
	if err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	// Output in shell format
	for key, value := range vars {
		fmt.Printf("export %s='%s'\n", key, value)
	}

	return nil
}

func init() {
	// Add flags
	envCmd.Flags().StringVarP(&writeFlag, "write", "w", "", "Set environment variable (format: KEY=VALUE)")
	envCmd.Flags().StringVarP(&unsetFlag, "unset", "u", "", "Remove environment variable")
	envCmd.Flags().BoolVar(&shellFlag, "shell", false, "Output environment variables in shell format")

	// Make flags mutually exclusive
	envCmd.MarkFlagsMutuallyExclusive("write", "unset", "shell")

	rootCmd.AddCommand(envCmd)
}

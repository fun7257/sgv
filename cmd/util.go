package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"golang.org/x/mod/semver"
)

// findGoModVersion searches for a go.mod file in the current directory
// and returns the Go version specified in it.
func findGoModVersion() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	goModPath := filepath.Join(currentDir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		content, err := os.ReadFile(goModPath)
		if err != nil {
			return "", fmt.Errorf("failed to read go.mod file: %w", err)
		}

		goRe := regexp.MustCompile(`^go\s+(\d+\.\d+(?:\.\d+)?)`)
		toolchainRe := regexp.MustCompile(`^toolchain\s+go(\d+\.\d+(?:\.\d+)?)`)

		var goVersion, toolchainVersion string

		lines := strings.SplitSeq(string(content), "\n")
		for line := range lines {
			if goMatches := goRe.FindStringSubmatch(line); len(goMatches) > 1 {
				goVersion = "go" + goMatches[1]
			}
			if toolchainMatches := toolchainRe.FindStringSubmatch(line); len(toolchainMatches) > 1 {
				toolchainVersion = "go" + toolchainMatches[1]
			}
		}

		if toolchainVersion != "" && (goVersion == "" || semver.Compare(normalizeGoVersion(toolchainVersion), normalizeGoVersion(goVersion)) > 0) {
			return normalizeGoModVersion(toolchainVersion), nil
		}

		if goVersion != "" {
			return normalizeGoModVersion(goVersion), nil
		}

		// If go.mod found but no go version specified, return empty string and nil error
		return "", nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("error checking for go.mod file: %w", err)
	}

	return "", nil // No go.mod found in current directory
}

// normalizeGoModVersion adjusts the Go version string based on the go.mod parsing rules.
// For Go 1.21 and later, "go1.21" is treated as "go1.21.0".
// For versions below 1.21, "go1.20" is treated as "go1.20".
func normalizeGoModVersion(v string) string {
	// Remove "go" prefix for easier parsing
	versionNum := strings.TrimPrefix(v, "go")

	// Check if it's a major.minor version (e.g., "1.21")
	if !strings.Contains(versionNum, ".") {
		return v // Should not happen for valid Go versions
	}

	parts := strings.Split(versionNum, ".")
	if len(parts) == 2 && semver.Compare("v"+versionNum, "v1.21") >= 0 { // Major.Minor version
		return "go" + versionNum + ".0"
	}
	return v
}

// normalizeGoVersion ensures the version string starts with 'v' for semver comparison.
func normalizeGoVersion(v string) string {
	if !strings.HasPrefix(v, "go") {
		v = "go" + v
	}
	if !strings.HasPrefix(v, "v") {
		return "v" + strings.TrimPrefix(v, "go")
	}
	return v
}

// isGoVersionCompatible checks if candidateVersion is greater than or equal to requiredVersion.
func isGoVersionCompatible(candidateVersion, requiredVersion string) bool {
	candidate := normalizeGoVersion(candidateVersion)
	required := normalizeGoVersion(requiredVersion)
	return semver.Compare(candidate, required) >= 0
}

// isGoVersionSupported checks if the given Go version is 1.13 or later.
func isGoVersionSupported(v string) bool {
	// Normalize to semver format (e.g., "go1.13" -> "v1.13")
	normalizedVersion := normalizeGoVersion(v)

	// Define the minimum supported version
	minSupportedVersion := normalizeGoVersion("go1.13")

	// Compare using semver.Compare
	return semver.Compare(normalizedVersion, minSupportedVersion) >= 0
}

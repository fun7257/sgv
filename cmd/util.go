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

		re := regexp.MustCompile(`^go\s+(\d+\.\d+(\.\d+)?)$`)
		lines := strings.SplitSeq(string(content), "\n")
		for line := range lines {
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				return "go" + matches[1], nil
			}
		}
		// If go.mod found but no go version specified, return empty string and nil error
		return "", nil
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("error checking for go.mod file: %w", err)
	}

	return "", nil // No go.mod found in current directory
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

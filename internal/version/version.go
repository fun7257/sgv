package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/fun7257/sgv/internal/config"
)

// GetLocalVersions reads the VERSIONS_DIR and returns a sorted list of installed version names.
func GetLocalVersions() ([]string, error) {
	files, err := os.ReadDir(config.VERSIONS_DIR)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions directory: %w", err)
	}

	var versions []string
	for _, file := range files {
		if file.IsDir() {
			versions = append(versions, file.Name())
		}
	}

	sort.Strings(versions)
	return versions, nil
}

// GetCurrentVersion reads the CURRENT_SYMLINK and returns the name of the active version.
func GetCurrentVersion() (string, error) {
	linkPath, err := os.Readlink(config.CURRENT_SYMLINK)
	if err != nil {
		return "", fmt.Errorf("failed to read current symlink: %w", err)
	}

	// The symlink points to a directory inside VERSIONS_DIR, so we need to get the base name.
	// We need to get the version name (e.g., "go1.17") from the symlink path.
	versionDir := filepath.Base(filepath.Dir(linkPath))
	return versionDir, nil
}

// GoVersion represents a Go version from the remote API.
type GoVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

// GetRemoteVersions fetches available Go versions from the official Go website.
func GetRemoteVersions() ([]GoVersion, error) {
	resp, err := http.Get("https://go.dev/dl/?mode=json")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch remote versions: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var versions []GoVersion
	if err := json.Unmarshal(body, &versions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	return versions, nil
}

// SwitchToVersion removes the existing CURRENT_SYMLINK and creates a new one.
func SwitchToVersion(version string) error {
	// Remove existing symlink if it exists
	if _, err := os.Lstat(config.CURRENT_SYMLINK); err == nil {
		if err := os.Remove(config.CURRENT_SYMLINK); err != nil {
			return fmt.Errorf("failed to remove existing symlink: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat current symlink: %w", err)
	}

	// Create new symlink
	targetPath := filepath.Join(config.VERSIONS_DIR, version, "go") // Symlink to the 'go' directory inside the version
	if err := os.Symlink(targetPath, config.CURRENT_SYMLINK); err != nil {
		return fmt.Errorf("failed to create new symlink: %w", err)
	}

	return nil
}

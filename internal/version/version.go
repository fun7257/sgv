package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"sort"

	"github.com/fun7257/sgv/internal/config"
)

var (
	sgvVersion = "dev" // Variable to hold SGV's own version
	sgvCommit  = "none"
	goVersion  = "unknown" // Variable to hold the Go version used to build SGV
)

// GetSGVVersion reads build info and returns SGV's version string.
func GetSGVVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" {
			sgvVersion = info.Main.Version
		}
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				sgvCommit = setting.Value
			}
		}

		goVersion = info.GoVersion
	}
	return fmt.Sprintf("%s (commit: %s, goVersion: %s)", sgvVersion, sgvCommit, goVersion)
}

// GetLocalVersions reads the VersionsDir and returns a sorted list of installed version names.
func GetLocalVersions() ([]string, error) {
	files, err := os.ReadDir(config.VersionsDir)
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

// GetCurrentVersion reads the CurrentSymlink and returns the name of the active version.
func GetCurrentVersion() (string, error) {
	linkPath, err := os.Readlink(config.CurrentSymlink)
	if err != nil {
		return "", fmt.Errorf("failed to read current symlink: %w", err)
	}

	// The symlink points to a directory inside VersionsDir, so we need to get the base name.
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

// GetLatestGoVersion fetches the latest stable Go version from the official Go website.
func GetLatestGoVersion() (string, error) {
	versions, err := GetRemoteVersions()
	if err != nil {
		return "", err
	}

	for _, v := range versions {
		if v.Stable {
			return v.Version, nil
		}
	}

	return "", fmt.Errorf("no stable Go version found")
}

// SwitchToVersion removes the existing CurrentSymlink and creates a new one.
func SwitchToVersion(version string) error {
	// Remove existing symlink if it exists
	if _, err := os.Lstat(config.CurrentSymlink); err == nil {
		if err := os.Remove(config.CurrentSymlink); err != nil {
			return fmt.Errorf("failed to remove existing symlink: %w", err)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat current symlink: %w", err)
	}

	// Create new symlink
	targetPath := filepath.Join(config.VersionsDir, version, "go") // Symlink to the 'go' directory inside the version
	if err := os.Symlink(targetPath, config.CurrentSymlink); err != nil {
		return fmt.Errorf("failed to create new symlink: %w", err)
	}

	return nil
}

// fetchAllGoVersions fetches the content of the Go downloads page and extracts all version numbers.
func FetchAllGoVersions() ([]string, error) {
	resp, err := http.Get("https://go.dev/dl/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Regex to find download links like "/dl/go1.22.4.src.tar.gz"
	re := regexp.MustCompile(`"/dl/(go[0-9]+\.[0-9]+(\.[0-9]+)?)\.src\.tar\.gz"`)
	matches := re.FindAllStringSubmatch(string(body), -1)

	versions := make(map[string]struct{})
	for _, match := range matches {
		if len(match) > 1 {
			versions[match[1]] = struct{}{}
		}
	}

	versionList := make([]string, 0, len(versions))
	for v := range versions {
		versionList = append(versionList, v)
	}

	return versionList, nil
}

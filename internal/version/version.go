package version

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/fun7257/sgv/internal/config"
	"github.com/samber/lo"
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

// GoVersion represents a simplified Go version for internal use.
type GoVersion struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
	OS      string `json:"os"`
	Arch    string `json:"arch"`
}

// GoVersionFile represents a file download for a specific Go version
type GoVersionFile struct {
	Filename string `json:"filename"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Version  string `json:"version"`
	SHA256   string `json:"sha256"`
	Size     int64  `json:"size"`
	Kind     string `json:"kind"`
}

// GoVersionResponse represents the complete API response structure
type GoVersionResponse struct {
	Version string          `json:"version"`
	Stable  bool            `json:"stable"`
	Files   []GoVersionFile `json:"files"`
}

// GetRemoteVersions fetches available Go versions from the official Go website, with a 10-minute file cache.
func GetRemoteVersions() ([]GoVersion, error) {
	cache := NewVersionCache()

	// Try to load from cache first
	if versions, err := cache.LoadFresh(); err == nil {
		return versions, nil
	}

	// Fetch from remote API
	versions, err := fetchRemoteVersions()
	if err != nil {
		// If remote fetch fails, try to return stale cache as fallback
		if staleVersions, cacheErr := cache.LoadStale(); cacheErr == nil {
			return staleVersions, nil
		}
		return nil, err
	}

	// Save to cache (best effort, don't fail on cache write errors)
	cache.Save(versions)

	return versions, nil
}

// fetchRemoteVersions fetches versions from the official Go website
func fetchRemoteVersions() ([]GoVersion, error) {
	url := config.DownloadURLPrefix + "?mode=json&include=all"

	// Create request with timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch remote versions: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the complete API response
	var apiResponse []GoVersionResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON response: %w", err)
	}

	// Convert to simplified GoVersion slice
	var versions []GoVersion
	for _, response := range apiResponse {
		// Skip empty entries (first element is often empty)
		if response.Version == "" {
			continue
		}

		// Create GoVersion entries for each file (OS/Arch combination)
		for _, file := range response.Files {
			// Skip source files (they don't have OS/Arch)
			if file.Kind == "source" {
				continue
			}

			// Skip Windows files as Windows is not supported
			if file.OS == "windows" {
				continue
			}

			versions = append(versions, GoVersion{
				Version: response.Version,
				Stable:  response.Stable,
				OS:      file.OS,
				Arch:    file.Arch,
			})
		}
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
	// Basic input validation: avoid path separators in version
	if strings.Contains(version, string(os.PathSeparator)) {
		return fmt.Errorf("invalid version: %q", version)
	}

	// Ensure target exists and looks like an installed Go distribution
	targetPath := filepath.Join(config.VersionsDir, version, "go") // Symlink to the 'go' directory inside the version
	if fi, err := os.Stat(targetPath); err != nil {
		return fmt.Errorf("version %s is not installed at %s: %w", version, targetPath, err)
	} else if !fi.IsDir() {
		return fmt.Errorf("version target is not a directory: %s", targetPath)
	}

	// Ensure the go binary exists under targetPath/bin/go
	goBin := filepath.Join(targetPath, "bin", "go")
	if _, err := os.Stat(goBin); err != nil {
		return fmt.Errorf("go binary not found for version %s at %s: %w", version, goBin, err)
	}

	// Check existing CurrentSymlink: if it exists, ensure it's a symlink; don't remove arbitrary files
	if info, err := os.Lstat(config.CurrentSymlink); err == nil {
		if info.Mode()&os.ModeSymlink == 0 {
			return fmt.Errorf("current path exists and is not a symlink: %s", config.CurrentSymlink)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat current symlink: %w", err)
	}

	// Create a temporary symlink in the same directory as CurrentSymlink, then atomically rename into place
	dir := filepath.Dir(config.CurrentSymlink)
	tmpName := fmt.Sprintf(".current.tmp.%d.%d", os.Getpid(), time.Now().UnixNano())
	tmpPath := filepath.Join(dir, tmpName)

	if err := os.Symlink(targetPath, tmpPath); err != nil {
		return fmt.Errorf("failed to create temporary symlink %s -> %s: %w", tmpPath, targetPath, err)
	}

	// Atomically replace (rename will overwrite existing target on Unix-like systems)
	if err := os.Rename(tmpPath, config.CurrentSymlink); err != nil {
		// best-effort cleanup
		_ = os.Remove(tmpPath)
		return fmt.Errorf("failed to atomically rename %s to %s: %w", tmpPath, config.CurrentSymlink, err)
	}

	return nil
}

// GetStableGoVersions fetches all stable Go versions from the official Go website.
func GetStableGoVersions() ([]GoVersion, error) {
	remoteVersions, err := GetRemoteVersions()
	if err != nil {
		return nil, err
	}

	// Filter out non-stable versions
	remoteVersions = lo.Filter(remoteVersions, func(item GoVersion, _ int) bool {
		return item.Stable
	})

	return remoteVersions, nil
}

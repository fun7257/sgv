package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// VersionCache handles caching of Go version data
type VersionCache struct {
	cacheFile     string
	cacheDuration time.Duration
}

// NewVersionCache creates a new version cache instance
func NewVersionCache() *VersionCache {
	return &VersionCache{
		cacheFile:     filepath.Join(os.TempDir(), "sgv-remote-versions-cache.json"),
		cacheDuration: time.Hour,
	}
}

// LoadFresh loads versions from cache if it exists and is fresh
func (c *VersionCache) LoadFresh() ([]GoVersion, error) {
	fi, err := os.Stat(c.cacheFile)
	if err != nil {
		return nil, fmt.Errorf("cache file not found: %w", err)
	}

	// Check if cache is fresh
	if time.Since(fi.ModTime()) > c.cacheDuration {
		return nil, fmt.Errorf("cache expired")
	}

	data, err := os.ReadFile(c.cacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var versions []GoVersion
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	return versions, nil
}

// LoadStale loads versions from cache even if expired (fallback)
func (c *VersionCache) LoadStale() ([]GoVersion, error) {
	data, err := os.ReadFile(c.cacheFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read stale cache: %w", err)
	}

	var versions []GoVersion
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal stale cache: %w", err)
	}

	return versions, nil
}

// Save saves versions to cache file (best effort)
func (c *VersionCache) Save(versions []GoVersion) {
	data, err := json.Marshal(versions)
	if err != nil {
		return // Ignore marshal errors
	}

	// Ensure cache directory exists
	if err := os.MkdirAll(filepath.Dir(c.cacheFile), 0755); err != nil {
		return // Ignore directory creation errors
	}

	// Write to temporary file first, then rename for atomic operation
	tempFile := c.cacheFile + ".tmp"
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		return // Ignore write errors
	}

	// Atomic rename
	if err := os.Rename(tempFile, c.cacheFile); err != nil {
		// Clean up temp file if rename fails
		_ = os.Remove(tempFile)
	}
}

// Clear removes the cache file
func (c *VersionCache) Clear() error {
	if err := os.Remove(c.cacheFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to clear cache: %w", err)
	}
	return nil
}

// IsExpired checks if the cache is expired
func (c *VersionCache) IsExpired() bool {
	fi, err := os.Stat(c.cacheFile)
	if err != nil {
		return true // Cache doesn't exist, consider it expired
	}
	return time.Since(fi.ModTime()) > c.cacheDuration
}

// SetCacheDuration allows customizing the cache duration
func (c *VersionCache) SetCacheDuration(duration time.Duration) {
	c.cacheDuration = duration
}

// GetCacheInfo returns cache file path and expiration status
func (c *VersionCache) GetCacheInfo() (string, bool, time.Time) {
	fi, err := os.Stat(c.cacheFile)
	if err != nil {
		return c.cacheFile, true, time.Time{}
	}
	return c.cacheFile, c.IsExpired(), fi.ModTime()
}

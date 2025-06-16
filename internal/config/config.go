package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	SGV_ROOT        string
	VERSIONS_DIR    string
	CURRENT_SYMLINK string
)

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user home directory: %v\n", err)
		os.Exit(1)
	}

	SGV_ROOT = filepath.Join(homeDir, ".sgv")
	VERSIONS_DIR = filepath.Join(SGV_ROOT, "versions")
	CURRENT_SYMLINK = filepath.Join(SGV_ROOT, "current")

	for _, dir := range []string{SGV_ROOT, VERSIONS_DIR} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
				os.Exit(1)
			}
		}
	}
}

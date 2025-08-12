package config

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	SgvRoot           string
	VersionsDir       string
	CurrentSymlink    string
	DownloadURLPrefix string
)

func Init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting user home directory: %v\n", err)
		os.Exit(1)
	}

	SgvRoot = filepath.Join(homeDir, ".sgv")
	VersionsDir = filepath.Join(SgvRoot, "versions")
	CurrentSymlink = filepath.Join(SgvRoot, "current")

	// Set DownloadURLPrefix from env or default
	DownloadURLPrefix = os.Getenv("SGV_DOWNLOAD_URL_PREFIX")
	if DownloadURLPrefix == "" {
		DownloadURLPrefix = "https://go.dev/dl/"
	}
	// Ensure ends with '/'
	if DownloadURLPrefix[len(DownloadURLPrefix)-1] != '/' {
		DownloadURLPrefix += "/"
	}

	for _, dir := range []string{SgvRoot, VersionsDir} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
				os.Exit(1)
			}
		}
	}
}

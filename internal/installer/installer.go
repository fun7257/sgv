package installer

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fun7257/sgv/internal/config"

	"github.com/schollz/progressbar/v3"
)

// Install downloads and installs the specified Go version.
func Install(version string) error {
	goOS := runtime.GOOS
	goARCH := runtime.GOARCH

	// Check if the current platform is supported
	if goOS == "windows" {
		return fmt.Errorf("Windows is not supported by sgv. This tool only works on macOS and Linux")
	}

	filename := fmt.Sprintf("%s.%s-%s.tar.gz", version, goOS, goARCH)
	downloadURL := fmt.Sprintf("%s%s", config.DownloadURLPrefix, filename)

	fmt.Printf("Downloading %s from %s\n", version, downloadURL)

	// Create the file to save the download
	outFilePath := filepath.Join(os.TempDir(), filename)
	out, err := os.Create(outFilePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outFilePath, err)
	}
	defer out.Close()

	// Download the file with a progress bar
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad HTTP status: %s", resp.Status)
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	if _, err := io.Copy(io.MultiWriter(out, bar), resp.Body); err != nil {
		return fmt.Errorf("failed to write download to file: %w", err)
	}
	out.Close()

	fmt.Printf("Extracting %s...\n", filename)

	// Extract the archive
	installPath := filepath.Join(config.VersionsDir, version)
	if err := extractTarGz(outFilePath, installPath); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	// Clean up the downloaded file
	defer os.Remove(outFilePath)

	return nil
}

// extractTarGz extracts a .tar.gz archive to the specified destination.
func extractTarGz(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // End of archive
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		targetPath := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}
		case tar.TypeReg:
			// Ensure parent directory exists
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for %s: %w", targetPath, err)
			}
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %w", targetPath, err)
			}
			// Close the file immediately after writing to release the file descriptor
			// This is crucial for handling a large number of files and preventing "too many open files" errors.
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close() // Close on error
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}
			outFile.Close()
			// Set file permissions
			if err := os.Chmod(targetPath, header.FileInfo().Mode()); err != nil {
				return fmt.Errorf("failed to set file permissions for %s: %w", targetPath, err)
			}
		case tar.TypeSymlink:
			// Handle symlinks
			if err := os.Symlink(header.Linkname, targetPath); err != nil {
				return fmt.Errorf("failed to create symlink %s: %w", targetPath, err)
			}
		default:
			return fmt.Errorf("unsupported tar entry type: %v in %s", header.Typeflag, header.Name)
		}
	}

	return nil
}

package version

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fun7257/sgv/internal/config"
)

func setupTemp(t *testing.T) string {
	t.Helper()
	tmp, err := os.MkdirTemp("", "sgv-test-")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(tmp) })
	return tmp
}

func TestSwitchToVersion(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		tmp := setupTemp(t)

		config.VersionsDir = filepath.Join(tmp, "versions")
		config.CurrentSymlink = filepath.Join(tmp, "current")

		// create installed version layout: <VersionsDir>/go1.13/go/bin/go
		targetBin := filepath.Join(config.VersionsDir, "go1.13", "go", "bin")
		if err := os.MkdirAll(targetBin, 0755); err != nil {
			t.Fatalf("mkdir target: %v", err)
		}
		goBin := filepath.Join(targetBin, "go")
		if err := os.WriteFile(goBin, []byte("#!/bin/sh\necho go"), 0755); err != nil {
			t.Fatalf("write go binary stub: %v", err)
		}

		// create an initial symlink to another version to ensure rename replaces it
		other := filepath.Join(config.VersionsDir, "go1.12", "go")
		if err := os.MkdirAll(other, 0755); err != nil {
			t.Fatalf("mkdir other: %v", err)
		}
		if err := os.Symlink(other, config.CurrentSymlink); err != nil {
			t.Fatalf("create initial symlink: %v", err)
		}

		if err := SwitchToVersion("go1.13"); err != nil {
			t.Fatalf("SwitchToVersion failed: %v", err)
		}

		// verify current symlink points to expected target
		linkTarget, err := os.Readlink(config.CurrentSymlink)
		if err != nil {
			t.Fatalf("readlink failed: %v", err)
		}
		expected := filepath.Join(config.VersionsDir, "go1.13", "go")
		if linkTarget != expected {
			t.Fatalf("unexpected symlink target: got %q want %q", linkTarget, expected)
		}

		// ensure no tmp files remain
		dir := filepath.Dir(config.CurrentSymlink)
		entries, _ := os.ReadDir(dir)
		for _, e := range entries {
			if strings.HasPrefix(e.Name(), ".current.tmp.") {
				t.Fatalf("temporary file left behind: %s", e.Name())
			}
		}
	})

	t.Run("target missing", func(t *testing.T) {
		tmp := setupTemp(t)

		config.VersionsDir = filepath.Join(tmp, "versions")
		config.CurrentSymlink = filepath.Join(tmp, "current")

		if err := SwitchToVersion("go9.9"); err == nil {
			t.Fatalf("expected error when switching to non-existent version")
		}

		if _, err := os.Lstat(config.CurrentSymlink); !os.IsNotExist(err) {
			t.Fatalf("expected no current symlink, got: %v", err)
		}
	})

	t.Run("current not symlink", func(t *testing.T) {
		tmp := setupTemp(t)

		config.VersionsDir = filepath.Join(tmp, "versions")
		config.CurrentSymlink = filepath.Join(tmp, "current")

		// prepare valid target
		targetBin := filepath.Join(config.VersionsDir, "go1.13", "go", "bin")
		if err := os.MkdirAll(targetBin, 0755); err != nil {
			t.Fatalf("mkdir target: %v", err)
		}
		goBin := filepath.Join(targetBin, "go")
		if err := os.WriteFile(goBin, []byte("#!/bin/sh\necho go"), 0755); err != nil {
			t.Fatalf("write go binary stub: %v", err)
		}

		// create a regular file at CurrentSymlink
		if err := os.WriteFile(config.CurrentSymlink, []byte("not a symlink"), 0644); err != nil {
			t.Fatalf("create regular file: %v", err)
		}

		if err := SwitchToVersion("go1.13"); err == nil {
			t.Fatalf("expected error when current path is not a symlink")
		}

		// ensure file content unchanged
		b, err := os.ReadFile(config.CurrentSymlink)
		if err != nil {
			t.Fatalf("read current file: %v", err)
		}
		if string(b) != "not a symlink" {
			t.Fatalf("current file was modified")
		}
	})
}

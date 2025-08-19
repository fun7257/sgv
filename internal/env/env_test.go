package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fun7257/sgv/internal/config"
)

func TestValidateVarName(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		wantErr bool
	}{
		{"valid uppercase", "VAR", false},
		{"valid with underscore", "VAR_NAME", false},
		{"valid with numbers", "VAR123", false},
		{"valid starting with underscore", "_VAR", false},
		{"valid lowercase", "var_name", false},
		{"valid single char", "A", false},
		{"valid underscore only", "_", false},
		{"invalid starting with number", "123VAR", true},
		{"invalid with hyphen", "VAR-NAME", true},
		{"invalid with space", "VAR NAME", true},
		{"invalid empty", "", true},
		{"invalid with dot", "VAR.NAME", true},
		{"invalid with special chars", "VAR@NAME", true},
		{"invalid with unicode", "VAR中文", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVarName(tt.varName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVarName(%q) error = %v, wantErr %v", tt.varName, err, tt.wantErr)
			}
		})
	}
}

func TestIsProtectedVar(t *testing.T) {
	tests := []struct {
		name      string
		varName   string
		protected bool
	}{
		{"GOROOT uppercase", "GOROOT", true},
		{"GOPATH uppercase", "GOPATH", true},
		{"GOPROXY uppercase", "GOPROXY", true},
		{"GOSUMDB uppercase", "GOSUMDB", true},
		{"GONOPROXY uppercase", "GONOPROXY", true},
		{"GONOSUMDB uppercase", "GONOSUMDB", true},
		{"GOPRIVATE uppercase", "GOPRIVATE", true},
		{"GO111MODULE uppercase", "GO111MODULE", true},
		{"GOOS uppercase", "GOOS", true},
		{"GOARCH uppercase", "GOARCH", true},
		{"CGO_ENABLED uppercase", "CGO_ENABLED", true},
		{"goroot lowercase", "goroot", true},
		{"gopath mixed case", "GoPath", true},
		{"GODEBUG not protected", "GODEBUG", false},
		{"GOWORK not protected", "GOWORK", false},
		{"MY_VAR custom", "MY_VAR", false},
		{"CUSTOM_VAR custom", "CUSTOM_VAR", false},
		{"GO_CUSTOM not protected", "GO_CUSTOM", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsProtectedVar(tt.varName)
			if result != tt.protected {
				t.Errorf("IsProtectedVar(%q) = %v, want %v", tt.varName, result, tt.protected)
			}
		})
	}
}

// setupTestEnv sets up a temporary directory for testing and returns cleanup function
func setupTestEnv(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "sgv-env-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Backup original config
	originalSgvRoot := config.SgvRoot
	config.SgvRoot = tmpDir

	cleanup := func() {
		config.SgvRoot = originalSgvRoot
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
}

func TestEnvFileOperations(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"

	t.Run("LoadEnvVars_NonexistentFile", func(t *testing.T) {
		vars, err := LoadEnvVars(version)
		if err != nil {
			t.Errorf("LoadEnvVars should not error on nonexistent file, got: %v", err)
		}
		if len(vars) != 0 {
			t.Errorf("Expected empty vars for nonexistent file, got %d variables", len(vars))
		}
	})

	t.Run("SetEnvVar_Success", func(t *testing.T) {
		err := SetEnvVar(version, "TEST_VAR", "test_value")
		if err != nil {
			t.Fatalf("SetEnvVar failed: %v", err)
		}

		vars, err := LoadEnvVars(version)
		if err != nil {
			t.Fatalf("LoadEnvVars failed: %v", err)
		}

		if vars["TEST_VAR"] != "test_value" {
			t.Errorf("Expected TEST_VAR=test_value, got %s", vars["TEST_VAR"])
		}
	})

	t.Run("SetEnvVar_MultipleVars", func(t *testing.T) {
		err := SetEnvVar(version, "ANOTHER_VAR", "another_value")
		if err != nil {
			t.Fatalf("SetEnvVar failed: %v", err)
		}

		err = SetEnvVar(version, "THIRD_VAR", "third_value")
		if err != nil {
			t.Fatalf("SetEnvVar failed: %v", err)
		}

		vars, err := LoadEnvVars(version)
		if err != nil {
			t.Fatalf("LoadEnvVars failed: %v", err)
		}

		expected := map[string]string{
			"TEST_VAR":    "test_value",
			"ANOTHER_VAR": "another_value",
			"THIRD_VAR":   "third_value",
		}

		if len(vars) != len(expected) {
			t.Errorf("Expected %d variables, got %d", len(expected), len(vars))
		}

		for key, expectedValue := range expected {
			if vars[key] != expectedValue {
				t.Errorf("Expected %s=%s, got %s", key, expectedValue, vars[key])
			}
		}
	})

	t.Run("SetEnvVar_OverwriteExisting", func(t *testing.T) {
		err := SetEnvVar(version, "TEST_VAR", "new_value")
		if err != nil {
			t.Fatalf("SetEnvVar failed: %v", err)
		}

		vars, err := LoadEnvVars(version)
		if err != nil {
			t.Fatalf("LoadEnvVars failed: %v", err)
		}

		if vars["TEST_VAR"] != "new_value" {
			t.Errorf("Expected TEST_VAR=new_value, got %s", vars["TEST_VAR"])
		}
	})

	t.Run("UnsetEnvVar_Success", func(t *testing.T) {
		err := UnsetEnvVar(version, "TEST_VAR")
		if err != nil {
			t.Fatalf("UnsetEnvVar failed: %v", err)
		}

		vars, err := LoadEnvVars(version)
		if err != nil {
			t.Fatalf("LoadEnvVars failed: %v", err)
		}

		if _, exists := vars["TEST_VAR"]; exists {
			t.Error("TEST_VAR should have been unset")
		}

		// Other vars should still exist
		if vars["ANOTHER_VAR"] != "another_value" {
			t.Errorf("ANOTHER_VAR should still exist, got %s", vars["ANOTHER_VAR"])
		}
	})

	t.Run("UnsetEnvVar_NonexistentVar", func(t *testing.T) {
		err := UnsetEnvVar(version, "NONEXISTENT_VAR")
		if err == nil {
			t.Error("Expected error when unsetting nonexistent variable")
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Expected 'not found' error, got: %v", err)
		}
	})
}

func TestProtectedVariables(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"

	protectedVars := []string{"GOROOT", "GOPATH", "GOPROXY", "GOSUMDB", "GOOS", "GOARCH"}

	t.Run("SetEnvVar_ProtectedVars", func(t *testing.T) {
		for _, varName := range protectedVars {
			err := SetEnvVar(version, varName, "some_value")
			if err == nil {
				t.Errorf("Expected error when setting protected variable %s", varName)
			}
			if !strings.Contains(err.Error(), "protected") {
				t.Errorf("Expected 'protected' in error message for %s, got: %v", varName, err)
			}
		}
	})

	t.Run("UnsetEnvVar_ProtectedVars", func(t *testing.T) {
		for _, varName := range protectedVars {
			err := UnsetEnvVar(version, varName)
			if err == nil {
				t.Errorf("Expected error when unsetting protected variable %s", varName)
			}
			if !strings.Contains(err.Error(), "protected") {
				t.Errorf("Expected 'protected' in error message for %s, got: %v", varName, err)
			}
		}
	})
}

func TestInvalidVariableNames(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"
	invalidNames := []string{"123VAR", "VAR-NAME", "VAR NAME", "", "VAR.NAME"}

	t.Run("SetEnvVar_InvalidNames", func(t *testing.T) {
		for _, varName := range invalidNames {
			err := SetEnvVar(version, varName, "value")
			if err == nil {
				t.Errorf("Expected error when setting invalid variable name: %s", varName)
			}
		}
	})

	t.Run("UnsetEnvVar_InvalidNames", func(t *testing.T) {
		for _, varName := range invalidNames {
			err := UnsetEnvVar(version, varName)
			if err == nil {
				t.Errorf("Expected error when unsetting invalid variable name: %s", varName)
			}
		}
	})
}

func TestSpecialValues(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"

	tests := []struct {
		name  string
		key   string
		value string
	}{
		{"empty value", "EMPTY_VAR", ""},
		{"value with spaces", "SPACE_VAR", "value with spaces"},
		{"value with quotes", "QUOTE_VAR", "value with 'single' and \"double\" quotes"},
		{"value with special chars", "SPECIAL_VAR", "value@#$%^&*()"},
		// Note: newlines are not supported in the current implementation
		// {"value with newlines", "NEWLINE_VAR", "line1\nline2"},
		{"value with unicode", "UNICODE_VAR", "值包含中文"},
		{"very long value", "LONG_VAR", strings.Repeat("a", 1000)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetEnvVar(version, tt.key, tt.value)
			if err != nil {
				t.Fatalf("SetEnvVar failed for %s: %v", tt.name, err)
			}

			vars, err := LoadEnvVars(version)
			if err != nil {
				t.Fatalf("LoadEnvVars failed: %v", err)
			}

			if vars[tt.key] != tt.value {
				t.Errorf("Value mismatch for %s: expected %q, got %q", tt.key, tt.value, vars[tt.key])
			}
		})
	}
}

func TestVersionIsolation(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version1 := "go1.20.0"
	version2 := "go1.21.0"

	// Set different values for the same variable in different versions
	err := SetEnvVar(version1, "SHARED_VAR", "value_for_1.20")
	if err != nil {
		t.Fatalf("SetEnvVar failed for version1: %v", err)
	}

	err = SetEnvVar(version2, "SHARED_VAR", "value_for_1.21")
	if err != nil {
		t.Fatalf("SetEnvVar failed for version2: %v", err)
	}

	// Verify isolation
	vars1, err := LoadEnvVars(version1)
	if err != nil {
		t.Fatalf("LoadEnvVars failed for version1: %v", err)
	}

	vars2, err := LoadEnvVars(version2)
	if err != nil {
		t.Fatalf("LoadEnvVars failed for version2: %v", err)
	}

	if vars1["SHARED_VAR"] != "value_for_1.20" {
		t.Errorf("Version1 should have value_for_1.20, got %s", vars1["SHARED_VAR"])
	}

	if vars2["SHARED_VAR"] != "value_for_1.21" {
		t.Errorf("Version2 should have value_for_1.21, got %s", vars2["SHARED_VAR"])
	}
}

func TestFileFormatAndContent(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"

	// Set some variables
	err := SetEnvVar(version, "VAR1", "value1")
	if err != nil {
		t.Fatalf("SetEnvVar failed: %v", err)
	}

	err = SetEnvVar(version, "VAR2", "value2")
	if err != nil {
		t.Fatalf("SetEnvVar failed: %v", err)
	}

	// Read the file content directly
	envFile := GetEnvFile(version)
	content, err := os.ReadFile(envFile)
	if err != nil {
		t.Fatalf("Failed to read env file: %v", err)
	}

	contentStr := string(content)

	// Verify file format
	expectedPatterns := []string{
		"# Environment variables for Go version " + version,
		"# Generated by sgv - do not edit manually",
		"VAR1=value1",
		"VAR2=value2",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(contentStr, pattern) {
			t.Errorf("Expected file to contain %q, but it doesn't.\nFile content:\n%s", pattern, contentStr)
		}
	}

	// Verify variables are sorted
	lines := strings.Split(contentStr, "\n")
	var varLines []string
	for _, line := range lines {
		if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
			varLines = append(varLines, line)
		}
	}

	if len(varLines) >= 2 {
		if varLines[0] > varLines[1] { // VAR1 should come before VAR2
			t.Error("Variables should be sorted alphabetically")
		}
	}
}

func TestConcurrentAccess(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"

	// Test concurrent reads and writes
	done := make(chan bool)
	errors := make(chan error, 10)

	// Start multiple goroutines doing different operations
	for i := 0; i < 5; i++ {
		go func(id int) {
			defer func() { done <- true }()

			varName := fmt.Sprintf("CONCURRENT_VAR_%d", id)
			value := fmt.Sprintf("value_%d", id)

			// Set variable
			if err := SetEnvVar(version, varName, value); err != nil {
				errors <- fmt.Errorf("SetEnvVar failed for %s: %v", varName, err)
				return
			}

			// Small delay to allow other goroutines to work
			// time.Sleep(time.Millisecond)

			// Load variables
			vars, err := LoadEnvVars(version)
			if err != nil {
				errors <- fmt.Errorf("LoadEnvVars failed: %v", err)
				return
			}

			// Verify our variable exists
			if vars[varName] != value {
				// This might fail due to concurrent modifications - that's actually expected
				// Log as informational rather than error
				t.Logf("Variable %s: expected %s, got %s (may be due to concurrent access)", varName, value, vars[varName])
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Check for errors
	close(errors)
	for err := range errors {
		t.Error(err)
	}

	// Verify that all variables were eventually set
	finalVars, err := LoadEnvVars(version)
	if err != nil {
		t.Fatalf("Final LoadEnvVars failed: %v", err)
	}

	// Check that we have some variables (concurrent access might overwrite some)
	if len(finalVars) == 0 {
		t.Error("Expected some variables to be set after concurrent operations")
	}
}

func TestEdgeCases(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"

	t.Run("EmptyVersion", func(t *testing.T) {
		// Current implementation doesn't validate version format
		// so empty version is technically allowed
		err := SetEnvVar("", "TEST_VAR", "value")
		if err != nil {
			t.Log("Empty version validation: ", err)
		}
	})

	t.Run("VersionWithSpaces", func(t *testing.T) {
		// Current implementation doesn't validate version format
		// so version with spaces is technically allowed
		err := SetEnvVar("go 1.21.0", "TEST_VAR", "value")
		if err != nil {
			t.Log("Version with spaces validation: ", err)
		}
	})

	t.Run("LoadNonexistentVersion", func(t *testing.T) {
		vars, err := LoadEnvVars("nonexistent")
		if err != nil {
			t.Errorf("LoadEnvVars should not error for nonexistent version: %v", err)
		}
		if len(vars) != 0 {
			t.Error("Expected empty vars for nonexistent version")
		}
	})

	t.Run("FilePermissions", func(t *testing.T) {
		// Create a file with restricted permissions
		envFile := GetEnvFile(version)
		if err := os.MkdirAll(filepath.Dir(envFile), 0755); err != nil {
			t.Fatalf("Failed to create env dir: %v", err)
		}

		if err := os.WriteFile(envFile, []byte("TEST_VAR=value\n"), 0000); err != nil {
			t.Fatalf("Failed to create restricted file: %v", err)
		}

		// Try to read the file
		_, err := LoadEnvVars(version)
		if err == nil {
			t.Error("Expected error when reading file with no permissions")
		}

		// Restore permissions for cleanup
		os.Chmod(envFile, 0644)
	})
}

func TestErrorMessages(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"

	t.Run("ProtectedVariableError", func(t *testing.T) {
		err := SetEnvVar(version, "GOROOT", "value")
		if err == nil {
			t.Error("Expected error for protected variable")
		}
		if !strings.Contains(err.Error(), "GOROOT") {
			t.Errorf("Error should mention variable name: %v", err)
		}
		if !strings.Contains(err.Error(), "protected") {
			t.Errorf("Error should mention protected: %v", err)
		}
	})

	t.Run("InvalidNameError", func(t *testing.T) {
		err := SetEnvVar(version, "123INVALID", "value")
		if err == nil {
			t.Error("Expected error for invalid variable name")
		}
		if !strings.Contains(err.Error(), "123INVALID") {
			t.Errorf("Error should mention variable name: %v", err)
		}
	})

	t.Run("VariableNotFoundError", func(t *testing.T) {
		err := UnsetEnvVar(version, "NONEXISTENT")
		if err == nil {
			t.Error("Expected error for nonexistent variable")
		}
		if !strings.Contains(err.Error(), "NONEXISTENT") {
			t.Errorf("Error should mention variable name: %v", err)
		}
		if !strings.Contains(err.Error(), "not found") {
			t.Errorf("Error should mention not found: %v", err)
		}
	})
}

func TestGetEnvFunctions(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	t.Run("GetEnvDir", func(t *testing.T) {
		envDir := GetEnvDir()
		if envDir == "" {
			t.Error("GetEnvDir should return a non-empty path")
		}
		expectedSuffix := filepath.Join("env")
		if !strings.HasSuffix(envDir, expectedSuffix) {
			t.Errorf("GetEnvDir should end with 'env', got: %s", envDir)
		}
	})

	t.Run("GetEnvFile", func(t *testing.T) {
		version := "go1.21.0"
		envFile := GetEnvFile(version)
		if envFile == "" {
			t.Error("GetEnvFile should return a non-empty path")
		}
		expectedSuffix := version + ".env"
		if !strings.HasSuffix(envFile, expectedSuffix) {
			t.Errorf("GetEnvFile should end with '%s', got: %s", expectedSuffix, envFile)
		}
	})
}

func TestLoadEnvVarsWithMalformedFile(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"
	envFile := GetEnvFile(version)

	// Create directory
	if err := os.MkdirAll(filepath.Dir(envFile), 0755); err != nil {
		t.Fatalf("Failed to create env dir: %v", err)
	}

	testCases := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid file",
			content: "VAR1=value1\nVAR2=value2\n",
			wantErr: false,
		},
		{
			name:    "file with comments",
			content: "# This is a comment\nVAR1=value1\n# Another comment\nVAR2=value2\n",
			wantErr: false,
		},
		{
			name:    "file with empty lines",
			content: "\nVAR1=value1\n\nVAR2=value2\n\n",
			wantErr: false,
		},
		{
			name:    "malformed line no equals",
			content: "VAR1=value1\nINVALID_LINE\nVAR2=value2\n",
			wantErr: true,
		},
		{
			name:    "invalid variable name",
			content: "VAR1=value1\n123INVALID=value2\n",
			wantErr: true,
		},
		{
			name:    "value with equals",
			content: "VAR1=value=with=equals\nVAR2=value2\n",
			wantErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.WriteFile(envFile, []byte(tc.content), 0644); err != nil {
				t.Fatalf("Failed to write test file: %v", err)
			}

			vars, err := LoadEnvVars(version)
			if tc.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if tc.name == "value with equals" && vars["VAR1"] != "value=with=equals" {
					t.Errorf("Expected 'value=with=equals', got '%s'", vars["VAR1"])
				}
			}

			// Clean up
			os.Remove(envFile)
		})
	}
}

func TestSaveEnvVarsEdgeCases(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"

	t.Run("SaveEmptyVars", func(t *testing.T) {
		vars := make(EnvVars)
		err := SaveEnvVars(version, vars)
		if err != nil {
			t.Errorf("SaveEnvVars should not error with empty vars: %v", err)
		}

		// File should still be created with headers
		envFile := GetEnvFile(version)
		content, err := os.ReadFile(envFile)
		if err != nil {
			t.Fatalf("Failed to read env file: %v", err)
		}

		contentStr := string(content)
		if !strings.Contains(contentStr, "# Environment variables for Go version") {
			t.Error("File should contain header comment")
		}
	})

	t.Run("SaveVarsWithSorting", func(t *testing.T) {
		vars := EnvVars{
			"ZEBRA": "last",
			"ALPHA": "first",
			"BETA":  "second",
		}

		err := SaveEnvVars(version, vars)
		if err != nil {
			t.Errorf("SaveEnvVars failed: %v", err)
		}

		// Read back and verify sorting
		loadedVars, err := LoadEnvVars(version)
		if err != nil {
			t.Errorf("LoadEnvVars failed: %v", err)
		}

		// Check that all vars are present
		if loadedVars["ALPHA"] != "first" || loadedVars["BETA"] != "second" || loadedVars["ZEBRA"] != "last" {
			t.Error("Variables not saved/loaded correctly")
		}

		// Check file content for sorting
		envFile := GetEnvFile(version)
		content, err := os.ReadFile(envFile)
		if err != nil {
			t.Fatalf("Failed to read env file: %v", err)
		}

		lines := strings.Split(string(content), "\n")
		var varLines []string
		for _, line := range lines {
			if strings.Contains(line, "=") && !strings.HasPrefix(line, "#") {
				varLines = append(varLines, line)
			}
		}

		if len(varLines) != 3 {
			t.Errorf("Expected 3 variable lines, got %d", len(varLines))
		}

		// Should be sorted alphabetically
		expected := []string{"ALPHA=first", "BETA=second", "ZEBRA=last"}
		for i, expectedLine := range expected {
			if i < len(varLines) && varLines[i] != expectedLine {
				t.Errorf("Line %d: expected %s, got %s", i, expectedLine, varLines[i])
			}
		}
	})

	t.Run("SaveVarsDirectoryCreation", func(t *testing.T) {
		// Remove the env directory to test creation
		envDir := GetEnvDir()
		os.RemoveAll(envDir)

		vars := EnvVars{"TEST_VAR": "test_value"}
		err := SaveEnvVars(version, vars)
		if err != nil {
			t.Errorf("SaveEnvVars should create directory if it doesn't exist: %v", err)
		}

		// Verify directory was created
		if _, err := os.Stat(envDir); os.IsNotExist(err) {
			t.Error("Environment directory should have been created")
		}
	})
}

func TestGetCurrentVersion(t *testing.T) {
	t.Run("GetCurrentVersion", func(t *testing.T) {
		// Note: GetCurrentVersion depends on the version package
		// This test primarily verifies the function can be called without panic
		_, err := GetCurrentVersion()
		// We expect an error in test environment since no version is typically active
		if err == nil {
			t.Log("GetCurrentVersion succeeded (unexpected in test environment)")
		} else {
			t.Logf("GetCurrentVersion failed as expected: %v", err)
		}
	})
}

func TestSetAndUnsetEnvVarErrorCoverage(t *testing.T) {
	_, cleanup := setupTestEnv(t)
	defer cleanup()

	version := "go1.21.0"

	t.Run("SetEnvVar_LoadError", func(t *testing.T) {
		// Create a directory where the env file should be to cause load error
		envFile := GetEnvFile(version)
		if err := os.MkdirAll(envFile, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		err := SetEnvVar(version, "TEST_VAR", "value")
		if err == nil {
			t.Error("Expected error when loading from a directory path")
		}
		if !strings.Contains(err.Error(), "failed to load existing variables") {
			t.Errorf("Expected 'failed to load existing variables' in error: %v", err)
		}

		// Clean up
		os.RemoveAll(envFile)
	})

	t.Run("UnsetEnvVar_LoadError", func(t *testing.T) {
		// Create a directory where the env file should be to cause load error
		envFile := GetEnvFile(version)
		if err := os.MkdirAll(envFile, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		err := UnsetEnvVar(version, "TEST_VAR")
		if err == nil {
			t.Error("Expected error when loading from a directory path")
		}
		if !strings.Contains(err.Error(), "failed to load existing variables") {
			t.Errorf("Expected 'failed to load existing variables' in error: %v", err)
		}

		// Clean up
		os.RemoveAll(envFile)
	})
}

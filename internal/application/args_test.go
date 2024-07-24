package application

import (
	"flag"
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigFileFromArgs(t *testing.T) {
	// Save the original command-line arguments
	origArgs := os.Args

	// Test cases
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "No config argument",
			args:     []string{"cmd"},
			expected: "",
		},
		{
			name:     "With config argument",
			args:     []string{"cmd", "-config=/path/to/config"},
			expected: "/path/to/config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the flag package to its default state
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ExitOnError)

			// Set the command-line arguments
			os.Args = tt.args

			// Call the function
			result := getConfigFileFromArgs()

			// Check the result
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}

	// Restore the original command-line arguments
	os.Args = origArgs
}

func TestTryConfigFile(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatalf("Failed to remove temporary directory: %v", err)
		}
	}(tempDir) // Clean up

	// Create a readable temporary file
	readableFile, err := os.CreateTemp(tempDir, "readable")
	if err != nil {
		t.Fatalf("Failed to create readable file: %v", err)
	}
	defer func(readableFile *os.File) {
		err := readableFile.Close()
		if err != nil {
			t.Fatalf("Failed to close readable file: %v", err)
		}
	}(readableFile)

	// Create a non-readable file
	nonReadableFile := filepath.Join(tempDir, "non_readable")
	err = os.WriteFile(nonReadableFile, []byte("data"), 0000) // No permissions
	if err != nil {
		t.Fatalf("Failed to create non-readable file: %v", err)
	}

	// Test cases
	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{
			name:     "Readable file",
			filePath: readableFile.Name(),
			expected: true,
		},
		{
			name:     "Non-readable file",
			filePath: nonReadableFile,
			expected: false,
		},
		{
			name:     "Non-existent file",
			filePath: filepath.Join(tempDir, "non_existent"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fullPath, result := tryConfigFile(tt.filePath)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
			if result && fullPath != tt.filePath {
				resolvedPath, _ := filepath.Abs(tt.filePath)
				if fullPath != resolvedPath {
					t.Errorf("expected full path %v, got %v", resolvedPath, fullPath)
				}
			}
		})
	}
}

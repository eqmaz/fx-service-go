package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

// TestExeDir tests the ExeDir function
func TestExeDir(t *testing.T) {
	dir, err := ExeDir()
	if err != nil {
		t.Errorf("ExeDir() returned an error: %v", err)
		return
	}

	// Check if the directory exists
	info, err := os.Stat(dir)
	if os.IsNotExist(err) {
		t.Errorf("ExeDir() returned a non-existing path: %s", dir)
	} else if err != nil {
		t.Errorf("Error checking the directory: %v", err)
	} else if !info.IsDir() {
		t.Errorf("ExeDir() returned a path that is not a directory: %s", dir)
	}
}

func TestIsExeInCwd(t *testing.T) {
	// Save the original working directory to restore it after the test
	originalCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting original working directory: %v", err)
	}
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			t.Errorf("Error changing directory to %s: %v", dir, err)
		}
	}(originalCwd)

	// Temporary directory for testing
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Error creating temporary directory: %v", err)
	}
	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			t.Errorf("Error removing temporary directory: %v", err)
		}
	}(tempDir)

	exePath, err := os.Executable()
	if err != nil {
		t.Fatalf("Error getting executable path: %v", err)
	}
	exeDir := filepath.Dir(exePath)

	t.Run("Executable in current working directory", func(t *testing.T) {
		err = os.Chdir(exeDir)
		if err != nil {
			t.Fatalf("Error changing directory to %s: %v", exeDir, err)
		}
		inCwd, err := IsExeInCwd()
		if err != nil {
			t.Fatalf("Error calling IsExeInCwd: %v", err)
		}
		if !inCwd {
			t.Errorf("Expected true, got false")
		}
	})

	t.Run("Executable not in current working directory", func(t *testing.T) {
		err = os.Chdir(tempDir)
		if err != nil {
			t.Fatalf("Error changing directory to %s: %v", tempDir, err)
		}
		notInCwd, err := IsExeInCwd()
		if err != nil {
			t.Fatalf("Error calling IsExeInCwd: %v", err)
		}
		if notInCwd {
			t.Errorf("Expected false, got true")
		}
	})
}

func TestIsReadableFile(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "readablefiletest")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer func() {
		err := os.Remove(tempFile.Name())
		if err != nil {
			t.Errorf("Error removing temporary file: %v", err)
		}
	}()

	// Test case: File exists and is readable
	t.Run("File exists and is readable", func(t *testing.T) {
		if !IsReadableFile(tempFile.Name()) {
			t.Errorf("Expected true, got false")
		}
	})

	// Test case: File does not exist
	t.Run("File does not exist", func(t *testing.T) {
		nonExistentFile := "nonexistentfile"
		if IsReadableFile(nonExistentFile) {
			t.Errorf("Expected false, got true")
		}
	})

	// Test case: File is not readable
	t.Run("File is not readable", func(t *testing.T) {
		// Create a temporary directory and a file inside it
		tempDir, err := os.MkdirTemp("", "unreadabledirtest")
		if err != nil {
			t.Fatalf("Error creating temporary directory: %v", err)
		}
		defer func() {
			err := os.RemoveAll(tempDir)
			if err != nil {
				t.Errorf("Error removing temporary directory: %v", err)
			}
		}()
		unreadableFile := filepath.Join(tempDir, "unreadablefile")
		_, err = os.Create(unreadableFile)
		if err != nil {
			t.Fatalf("Error creating temporary file: %v", err)
		}

		// Change the permissions of the file to make it unreadable
		err = os.Chmod(unreadableFile, 0300)
		if err != nil {
			t.Fatalf("Error changing file permissions: %v", err)
		}

		if IsReadableFile(unreadableFile) {
			t.Errorf("Expected false, got true")
		}
	})
}

func TestGetType(t *testing.T) {
	tests := []struct {
		name     string
		variable interface{}
		want     string
	}{
		{
			name:     "String",
			variable: "Hello",
			want:     "string",
		},
		{
			name:     "Int",
			variable: 42,
			want:     "int",
		},
		{
			name:     "Float",
			variable: 3.14,
			want:     "float64",
		},
		{
			name:     "Bool",
			variable: true,
			want:     "bool",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetType(tt.variable); got != tt.want {
				t.Errorf("GetType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRound2(t *testing.T) {
	tests := []struct {
		name      string
		num       float64
		precision int
		want      float64
	}{
		{
			name:      "Round to 2 decimal places",
			num:       3.14159,
			precision: 2,
			want:      3.14,
		},
		{
			name:      "Round to 0 decimal places",
			num:       3.14159,
			precision: 0,
			want:      3,
		},
		{
			name:      "Round to 4 decimal places",
			num:       3.14159,
			precision: 4,
			want:      3.1416,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Round(tt.num, tt.precision); got != tt.want {
				t.Errorf("Round() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberFormat(t *testing.T) {
	tests := []struct {
		name         string
		number       float64
		decimals     int
		decPoint     string
		thousandsSep string
		want         string
	}{
		{
			name:         "Format with 2 decimal places",
			number:       1234567.89,
			decimals:     2,
			decPoint:     ".",
			thousandsSep: ",",
			want:         "1,234,567.89",
		},
		{
			name:         "Format with 0 decimal places",
			number:       1234567.89,
			decimals:     0,
			decPoint:     ".",
			thousandsSep: ",",
			want:         "1,234,568",
		},
		{
			name:         "Format with 3 decimal places",
			number:       1234567.89,
			decimals:     3,
			decPoint:     ".",
			thousandsSep: ",",
			want:         "1,234,567.890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumberFormat(tt.number, tt.decimals, tt.decPoint, tt.thousandsSep); got != tt.want {
				t.Errorf("NumberFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}

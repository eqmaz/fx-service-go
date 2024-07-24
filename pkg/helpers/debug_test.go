package helpers

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCaptureOutput tests the CaptureOutput function
func TestCaptureOutput(t *testing.T) {
	tests := []struct {
		name     string
		function func()
		expected string
	}{
		{"PrintfString", func() { fmt.Printf("hello world") }, "hello world"},
		{"PrintlnString", func() { fmt.Println("hello world") }, "hello world\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := CaptureOutput(test.function)
			if output != test.expected {
				t.Errorf("CaptureOutput() = %v, want %v", output, test.expected)
			}
		})
	}
}

// TestDump tests the Dump function
func TestDump(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"String", "hello", "[string]\n\"hello\"\n"},
		{"Int", 123, "[int]\n123\n"},
		{"Struct", struct{ Name string }{"Go"}, "[struct { Name string }]\n{\n  \"Name\": \"Go\"\n}\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := CaptureOutput(func() { Dump(test.input) })
			if !strings.Contains(output, test.expected) {
				t.Errorf("Dump() = %v, want %v", output, test.expected)
			}
		})
	}
}

// TestDd tests the Dd function
func TestDd(t *testing.T) {
	if os.Getenv("TEST_DD") == "1" {
		Dd("test")
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestDd")
	cmd.Env = append(os.Environ(), "TEST_DD=1")
	err := cmd.Run()

	var e *exec.ExitError
	if errors.As(err, &e) && !e.Success() {
		return
	}
	t.Fatalf("expected Dd to exit with status code 1, but it did not")
}

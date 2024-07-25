package console

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"testing"
)

// Helper function to capture the output of fmt.Printf
func captureOutput(f func()) string {
	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	err := w.Close()
	if err != nil {
		return ""
	}
	os.Stdout = old

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r)
	if err != nil {
		return ""
	}
	return buf.String()
}

// Test for the send function
func TestSend(t *testing.T) {
	// Test case when suspended is false
	Resume()
	output := captureOutput(func() {
		send("Hello, send!")
	})
	expected := "Hello, send!\n"
	strippedOutput := stripTimestamp(output)
	if strippedOutput != expected {
		t.Errorf("Expected %q but got %q", expected, strippedOutput)
	}

	// Test case when suspended is true
	Suspend()
	output = captureOutput(func() {
		send("This should not be printed.")
	})
	if output != "" {
		t.Errorf("Expected no output but got %q", output)
	}

	Resume()
}

// Test for the colorize function
func TestColorize(t *testing.T) {
	text := "Hello"
	color := ColorGrey

	// Capture the current value of isATTY and restore it after the test
	originalIsATTY := isATTY

	// Test when isATTY is false
	isATTY = false
	expected := text
	result := colorize(text, color)
	if result != expected {
		t.Errorf("colorize(%q, %q) with isATTY=false - got result '%q'; wanted '%q'", text, color, result, expected)
	}

	// Test when isATTY is true
	isATTY = true
	expected = fmt.Sprintf("%s%s%s", color, text, ColorReset)
	result = colorize(text, color)
	if result != expected {
		t.Errorf("colorize(%q, %q) with isATTY=true - got result '%q'; wanted '%q'", text, color, result, expected)
	}

	// Restore the original value of isATTY
	isATTY = originalIsATTY
}

// Test for the toStr function
func TestToStr(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{123, "123"},
		{"test", "test"},
		{3.14, "3.14"},
		{true, "true"},
	}

	for _, test := range tests {
		result := toStr(test.input)
		if result != test.expected {
			t.Errorf("toStr(%v) = %q; want %q", test.input, result, test.expected)
		}
	}
}

// Test for the outTimestamp function
func TestOutTimestamp(t *testing.T) {
	output := outTimestamp()

	// Regular expression to match the expected timestamp format with milliseconds
	regex := `^\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}\]$`
	matched, _ := regexp.MatchString(regex, output)
	if !matched {
		t.Errorf("outTimestamp() = %q; does not match the expected format", output)
	}
}

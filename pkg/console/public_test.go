package console

import (
	"fmt"
	"regexp"
	"testing"
)

// Helper function to strip timestamp for comparison
func stripTimestamp(output string) string {
	re := regexp.MustCompile(`\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}\] `)
	return re.ReplaceAllString(output, "")
}

// Test for the Out function
func TestOut(t *testing.T) {
	testCases := [][]interface{}{
		{"Hello, world!"},
		{123, " ", 456},
		{"Pi is", 3.14},
	}

	for _, tc := range testCases {
		output := captureOutput(func() {
			Out(tc...)
		})

		// Regular expression to match the expected timestamp format with milliseconds
		regex := `^\[\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}\]`
		matched, _ := regexp.MatchString(regex, output)
		if !matched {
			t.Errorf("Out(%v) = %q; timestamp does not match the expected format", tc, output)
		}

		// Check if the output contains the string representations of the test case values
		for _, v := range tc {
			expectedStr := toStr(v)
			if !regexp.MustCompile(regexp.QuoteMeta(expectedStr)).MatchString(output) {
				t.Errorf("Out(%v) = %q; does not contain %q", tc, output, expectedStr)
			}
		}
	}
}

// Test for the Outf function
func TestOutf(t *testing.T) {
	format := "Hello, %s!"
	arg := "world"
	expected := fmt.Sprintf(format, arg) + "\n"

	output := captureOutput(func() {
		Outf(format, arg)
	})
	strippedOutput := stripTimestamp(output)

	if strippedOutput != expected {
		t.Errorf("Expected %q but got %q", expected, strippedOutput)
	}
}

// Test for the Success function
func TestSuccess(t *testing.T) {
	testCases := [][]interface{}{
		{"Operation completed successfully."},
		{"Task", 123, "completed."},
	}

	for _, testCase := range testCases {
		output := captureOutput(func() {
			Success(testCase...)
		})
		expected := "✔ " + fmt.Sprintln(testCase...)
		strippedOutput := stripTimestamp(output)
		if strippedOutput != expected {
			t.Errorf("Expected %q but got %q", expected, strippedOutput)
		}
	}
}

// Test for the Successf function
func TestSuccessf(t *testing.T) {
	format := "Operation %s successfully."
	arg := "completed"
	expected := "✔ " + fmt.Sprintf(format, arg) + "\n"

	output := captureOutput(func() {
		Successf(format, arg)
	})
	strippedOutput := stripTimestamp(output)

	if strippedOutput != expected {
		t.Errorf("Expected %q but got %q", expected, strippedOutput)
	}
}

// Test for the Warn function
func TestWarn(t *testing.T) {
	testCases := [][]interface{}{
		{"This is a warning."},
		{"Potential issue", 123},
	}

	for _, testCase := range testCases {
		output := captureOutput(func() {
			Warn(testCase...)
		})
		expected := "w " + fmt.Sprintln(testCase...)
		strippedOutput := stripTimestamp(output)
		if strippedOutput != expected {
			t.Errorf("Expected %q but got %q", expected, strippedOutput)
		}
	}
}

// Test for the Warnf function
func TestWarnf(t *testing.T) {
	format := "Warning: %s detected."
	arg := "issue"
	expected := "w " + fmt.Sprintf(format, arg) + "\n"

	output := captureOutput(func() {
		Warnf(format, arg)
	})
	strippedOutput := stripTimestamp(output)

	if strippedOutput != expected {
		t.Errorf("Expected %q but got %q", expected, strippedOutput)
	}
}

// Test for the Info function
func TestInfo(t *testing.T) {
	testCases := [][]interface{}{
		{"This is an info message."},
		{"Information", 123},
	}

	for _, testCase := range testCases {
		output := captureOutput(func() {
			Info(testCase...)
		})
		expected := "ℹ " + fmt.Sprintln(testCase...)
		strippedOutput := stripTimestamp(output)
		if strippedOutput != expected {
			t.Errorf("Expected %q but got %q", expected, strippedOutput)
		}
	}
}

// Test for the Infof function
func TestInfof(t *testing.T) {
	format := "Information: %s."
	arg := "details"
	expected := "ℹ " + fmt.Sprintf(format, arg) + "\n"

	output := captureOutput(func() {
		Infof(format, arg)
	})
	strippedOutput := stripTimestamp(output)

	if strippedOutput != expected {
		t.Errorf("Expected %q but got %q", expected, strippedOutput)
	}
}

// Test for the Suspend and Resume functions
func TestSuspendAndResume(t *testing.T) {
	Suspend()
	output := captureOutput(func() {
		Out("This should not be printed.")
	})
	if output != "" {
		t.Errorf("Expected no output but got %q", output)
	}

	Resume()
	output = captureOutput(func() {
		Out("This should be printed.")
	})
	expected := "This should be printed.\n"
	strippedOutput := stripTimestamp(output)
	if strippedOutput != expected {
		t.Errorf("Expected %q but got %q", expected, strippedOutput)
	}
}

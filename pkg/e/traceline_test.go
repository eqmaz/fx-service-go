package e

import (
	"testing"
)

func TestTraceLine_String(t *testing.T) {
	tests := []struct {
		name      string
		traceLine TraceLine
		expected  string
	}{
		{"basic case", TraceLine{File: "file.go", Line: 10, Func: "FuncName"}, "file.go:10 Func FuncName()"},
		{"empty values", TraceLine{File: "", Line: 0, Func: ""}, ":0 Func ()"},
		{"long function name", TraceLine{File: "file.go", Line: 123456, Func: "AReallyLongFunctionName"}, "file.go:123456 Func AReallyLongFunctionName()"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.traceLine.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

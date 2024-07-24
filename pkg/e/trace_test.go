package e

import (
	"testing"
)

func TestTrace_Top(t *testing.T) {
	tests := []struct {
		name     string
		trace    Trace
		expected *TraceLine
	}{
		{"empty trace", Trace{Lines: []TraceLine{}}, nil},
		{"single line", Trace{Lines: []TraceLine{{File: "file.go", Line: 1, Func: "Func1"}}}, &TraceLine{File: "file.go", Line: 1, Func: "Func1"}},
		{"multiple lines", Trace{Lines: []TraceLine{{File: "file1.go", Line: 1, Func: "Func1"}, {File: "file2.go", Line: 2, Func: "Func2"}}}, &TraceLine{File: "file1.go", Line: 1, Func: "Func1"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.trace.Top()
			if result == nil && tt.expected != nil {
				t.Errorf("expected %v, got nil", tt.expected)
			} else if result != nil && (*result != *tt.expected) {
				t.Errorf("expected %v, got %v", *tt.expected, *result)
			}
		})
	}
}

func TestTrace_Last(t *testing.T) {
	tests := []struct {
		name     string
		trace    Trace
		expected *TraceLine
	}{
		{"empty trace", Trace{Lines: []TraceLine{}}, nil},
		{"single line", Trace{Lines: []TraceLine{{File: "file.go", Line: 1, Func: "Func1"}}}, &TraceLine{File: "file.go", Line: 1, Func: "Func1"}},
		{"multiple lines", Trace{Lines: []TraceLine{{File: "file1.go", Line: 1, Func: "Func1"}, {File: "file2.go", Line: 2, Func: "Func2"}}}, &TraceLine{File: "file2.go", Line: 2, Func: "Func2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.trace.Last()
			if result == nil && tt.expected != nil {
				t.Errorf("expected %v, got nil", tt.expected)
			} else if result != nil && (*result != *tt.expected) {
				t.Errorf("expected %v, got %v", *tt.expected, *result)
			}
		})
	}
}

func TestTrace_Item(t *testing.T) {
	tests := []struct {
		name     string
		trace    Trace
		index    int
		expected *TraceLine
	}{
		{"negative index", Trace{Lines: []TraceLine{{File: "file.go", Line: 1, Func: "Func1"}}}, -1, nil},
		{"index out of bounds", Trace{Lines: []TraceLine{{File: "file.go", Line: 1, Func: "Func1"}}}, 1, nil},
		{"index 0 on single line", Trace{Lines: []TraceLine{{File: "file.go", Line: 1, Func: "Func1"}}}, 0, &TraceLine{File: "file.go", Line: 1, Func: "Func1"}},
		{"index 0 on multiple lines", Trace{Lines: []TraceLine{{File: "file1.go", Line: 1, Func: "Func1"}, {File: "file2.go", Line: 2, Func: "Func2"}}}, 0, &TraceLine{File: "file1.go", Line: 1, Func: "Func1"}},
		{"index 1 on multiple lines", Trace{Lines: []TraceLine{{File: "file1.go", Line: 1, Func: "Func1"}, {File: "file2.go", Line: 2, Func: "Func2"}}}, 1, &TraceLine{File: "file2.go", Line: 2, Func: "Func2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.trace.Item(tt.index)
			if result == nil && tt.expected != nil {
				t.Errorf("expected %v, got nil", tt.expected)
			} else if result != nil && (*result != *tt.expected) {
				t.Errorf("expected %v, got %v", *tt.expected, *result)
			}
		})
	}
}

func TestTrace_Count(t *testing.T) {
	tests := []struct {
		name     string
		trace    Trace
		expected int
	}{
		{"empty trace", Trace{Lines: []TraceLine{}}, 0},
		{"single line", Trace{Lines: []TraceLine{{File: "file.go", Line: 1, Func: "Func1"}}}, 1},
		{"multiple lines", Trace{Lines: []TraceLine{{File: "file1.go", Line: 1, Func: "Func1"}, {File: "file2.go", Line: 2, Func: "Func2"}}}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.trace.Count()
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

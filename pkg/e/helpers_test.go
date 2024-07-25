package e

import (
	"testing"
)

func Test_trimFilePath(t *testing.T) {
	tp := GetFilePathTrimPoint()

	SetFilePathTrimPoint("src/")
	tests := []struct {
		name     string
		fullPath string
		expected string
	}{
		{"trimmed path", "src/github.com/project/file.go", "github.com/project/file.go"},
		{"no trim needed", "github.com/project/file.go", "github.com/project/file.go"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := trimFilePath(tt.fullPath)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}

	// Put things back the way we found them
	SetFilePathTrimPoint(tp)
}

func Test_trimFunctionName(t *testing.T) {
	tests := []struct {
		input          string
		expectedPkg    string
		expectedStruct string
		expectedFunc   string
	}{
		// Struct method case
		{"github.com/user/project/package.(*StructName).MethodName",
			"package", "StructName", "MethodName"},
		// Lambda function case
		{"github.com/user/project/package.funcName.func1",
			"package", "", "funcName.func1"},
		// Nested function case
		{"github.com/user/project/package.funcName.func1.1",
			"package", "", "funcName.func1.1"},
		// Simple function case
		{"github.com/user/project/package.FunctionName",
			"package", "", "FunctionName"},
		// Nested packages with struct
		{"github.com/user/project/subpackage.(*StructName).MethodName",
			"subpackage", "StructName", "MethodName"},
		// Nested packages without struct
		{"github.com/user/project/subpackage.funcName",
			"subpackage", "", "funcName"},
		// Function with multiple dots
		{"github.com/user/project/package.func.with.dots",
			"package", "", "func.with.dots"},
	}

	for _, test := range tests {
		pkg, strct, fn := trimFunctionName(test.input)
		if pkg != test.expectedPkg {
			t.Errorf("For input '%s', expected package '%s', but got '%s'", test.input, test.expectedPkg, pkg)
		}
		if strct != test.expectedStruct {
			t.Errorf("For input '%s', expected struct '%s', but got '%s'", test.input, test.expectedStruct, strct)
		}
		if fn != test.expectedFunc {
			t.Errorf("For input '%s', expected function '%s', but got '%s'", test.input, test.expectedFunc, fn)
		}
	}
}

func Test_captureBacktrace(t *testing.T) {
	_, trace := captureBacktrace()
	if len(trace.Lines) == 0 {
		t.Error("expected non-empty trace lines")
	}
	firstLine := trace.Lines[0]
	if firstLine.File == "" {
		t.Error("expected non-empty file")
	}
	if firstLine.Func == "" {
		t.Error("expected non-empty function")
	}
	if firstLine.Line == 0 {
		t.Error("expected non-zero line number")
	}
}

func Test_makeException(t *testing.T) {
	trace := Trace{Lines: []TraceLine{{File: "file.go", Line: 10, Func: "FuncName"}}}
	line := 10
	code := "ERROR_CODE"
	msg := "An error occurred"
	ts := "trace string"
	file := "file.go"
	function := "FuncName"

	ex := makeException(trace, code, msg, ts)
	if ex == nil {
		t.Fatal("expected non-nil exception")
	}
	if ex.code != code {
		t.Errorf("expected code %q, got %q", code, ex.code)
	}
	if ex.message != msg {
		t.Errorf("expected message %q, got %q", msg, ex.message)
	}
	if ex.traceStr != ts {
		t.Errorf("expected trace string %q, got %q", ts, ex.traceStr)
	}
	if ex.origin.File != file {
		t.Errorf("expected file %q, got %q", file, ex.origin.File)
	}
	if ex.origin.Func != function {
		t.Errorf("expected function %q, got %q", function, ex.origin.Func)
	}
	if ex.origin.Line != uint(line) {
		t.Errorf("expected line %d, got %d", line, ex.origin.Line)
	}
}

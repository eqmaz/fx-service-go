package e

import (
	"strings"
	"testing"
)

func TestThrowBasic(t *testing.T) {
	ex := Throw("123", "test message")

	if ex.GetCode() != "123" {
		t.Errorf("expected code 123, got %v", ex.GetCode())
	}
	if ex.GetMessage() != "test message" {
		t.Errorf("expected message 'test message', got %s", ex.GetMessage())
	}
	if ex.GetPrevious() != nil {
		t.Errorf("expected previous error to be nil, got %v", ex.GetPrevious())
	}
	if ex.GetLine() == 0 {
		t.Errorf("expected line number to be greater than 0, got %d", ex.GetLine())
	}
	if ex.GetFile() == "" {
		t.Errorf("expected file name to be non-empty, got '%s'", ex.GetFile())
	}
	if ex.GetFunction() == "" {
		t.Errorf("expected function name to be non-empty, got '%s'", ex.GetFunction())
	}
}

func TestThrowWithPrevious(t *testing.T) {
	prev := Throw("456", "previous error")
	ex := Throw("789", "current error").SetPrevious(prev)

	if ex.GetCode() != "789" {
		t.Errorf("expected code 789, got %v", ex.GetCode())
	}
	if ex.GetMessage() != "current error" {
		t.Errorf("expected message 'current error', got %s", ex.GetMessage())
	}
	if ex.GetPrevious() == nil {
		t.Errorf("expected previous error to be non-nil")
	} else {
		prevEx := ex.GetPrevious()
		if prevEx == nil {
			t.Errorf("expected previous error to be of type *Exception, got nil")
			return
		}
		if prevEx.GetCode() != "456" {
			t.Errorf("expected previous code 456, got %v", prevEx.GetCode())
		}
		if prevEx.GetMessage() != "previous error" {
			t.Errorf("expected previous message 'previous error', got %s", prevEx.GetMessage())
		}
	}
}

func TestThrowBacktrace(t *testing.T) {
	ex := Throw("123", "test message with backtrace")

	trace := ex.GetTrace()
	if trace == nil {
		t.Errorf("expected backtrace to be non-empty")
		return
	}

	if trace.Count() == 0 {
		t.Errorf("expected backtrace to be non-empty")
	}
}

func TestThrowLine(t *testing.T) {
	ex := Throw("123", "test line")

	if ex.GetLine() == 0 {
		t.Errorf("expected line number to be greater than 0, got %d", ex.GetLine())
	}

	ex = Throw("123", "test line and file")
	expectedLine := 79 // NOTE - adjust if this file changes
	if ex.GetLine() != uint(expectedLine) {
		t.Errorf("expected line %d, got %d", expectedLine, ex.GetLine())
	}
}

func TestFile(t *testing.T) {
	ex := Throw("123", "test line and file")

	expectedFile := "/pkg/e/exception_test.go"
	// if not string contains
	if strings.Contains(ex.GetFile(), expectedFile) == false {
		t.Errorf("expected file %s, got %s", expectedFile, ex.GetFile())
	}
}

func TestThrowFunction(t *testing.T) {
	ex := Throw("123", "test function")

	expectedFunction := "TestThrowFunction"
	expectedPackage := "e"

	if ex.GetFunction() != expectedFunction {
		t.Errorf("expected function: '%s', got: '%s'", expectedFunction, ex.GetFunction())
	}
	if ex.GetPackage() != expectedPackage {
		t.Errorf("expected package: '%s', got: '%s'", expectedPackage, ex.GetPackage())
	}
}

func TestError(t *testing.T) {
	ex := Throw("123", "test error")

	expectedError := "Code: 123, Message: test error"

	if ex.Error() != expectedError {
		t.Errorf("expected error %s, got %s", expectedError, ex.Error())
	}
}

func TestJSON(t *testing.T) {
	SetFilePathTrimPoint("/fx-service")

	ex := Throw("abc", "test json")

	result, err := ex.JSON()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	expectedJson := `{"code":"abc","message":"test json","origin":{"File":"/pkg/e/exception_test.go","Line":123,"Package":"e","Struct":"","Func":"TestJSON"},"trace":[{"File":"/pkg/e/exception_test.go","Line":123,"Package":"e","Struct":"","Func":"TestJSON"},{"File":"/usr/local/go/src/testing/testing.go","Line":1689,"Package":"testing","Struct":"","Func":"tRunner"},{"File":"/usr/local/go/src/runtime/asm_amd64.s","Line":1695,"Package":"runtime","Struct":"","Func":"goexit"}]}`
	if result != expectedJson {
		t.Errorf("\nExpected JSON: %s \nGot: %s", expectedJson, result)
	}
}

func TestThrowChainedBacktrace(t *testing.T) {
	prev := Throw("456", "previous error with backtrace")
	ex := Throw("789", "current error with chained backtrace").SetPrevious(prev)

	if trace := ex.GetTrace(); trace == nil {
		t.Errorf("expected backtrace to be non-empty")
		return
	}
	if ex.GetPrevious() == nil {
		t.Errorf("expected previous error to be non-nil")
	} else {
		prevEx := ex.GetPrevious()
		if prevEx == nil {
			t.Errorf("expected previous error to be of type *Exception, got nil")
			return
		}
	}
}

package e

import (
	"encoding/json"
	"errors"
	"fmt"
	"fx-service/pkg/console"
	"github.com/mattn/go-isatty"
	"os"
	"sync"
)

// Adjust this to your project name or desired starting path
// It will ensure the file paths are trimmed to the relevant parts
// todo - make this a configuration option
var filePathTrimBefore = ""
var filePathTrimMutex = &sync.Mutex{}

// SetFilePathTrimPoint sets the pattern (needle) before which the file path will be trimmed.
// This is useful for trimming the file path to the relevant parts.
// It avoids having overly long and repetitive file paths in the exception trace.
// It's recommended to set this to the root of your project, or module name.
func SetFilePathTrimPoint(needle string) {
	filePathTrimMutex.Lock()
	filePathTrimBefore = needle
	filePathTrimMutex.Unlock()
}

// GetFilePathTrimPoint returns the current file path trim point.
func GetFilePathTrimPoint() string {
	return filePathTrimBefore
}

// Exception struct to represent an error with a code, message, previous exception, and backtrace.
type Exception struct {
	code     string
	message  string
	origin   *TraceLine
	previous *Exception
	traceStr string
	trace    *Trace
	fields   Fields
	args     []interface{}
}

// Throw creates the Exception struct.
// code should be a unique identifier for the error, message is a human-readable error message.
func Throw(code string, message string) *Exception {
	traceString, trace := captureBacktrace()
	return makeException(trace, code, message, traceString)
}

// Throwf creates the Exception struct with a formatted message.
func Throwf(code string, format string, args ...interface{}) *Exception {
	msg := fmt.Sprintf(format, args...)
	traceString, trace := captureBacktrace()
	return makeException(trace, code, msg, traceString)
}

// FromCode creates an Exception from a pre-catalogued error code.
// If the code is not found, a generic error will be returned.
// If there are arguments, an attempt will be made to format them into the message.
func FromCode(code string, args ...interface{}) *Exception {
	// See if the error code exists
	msg, ok := catalogue[code]
	if !ok {
		return Throwf("", "Unknown error code '%s'. Ensure error is catalogued. ", code)
	}

	// If there are arguments, format the message
	if args != nil {
		msg = fmt.Sprintf(msg, args...)
	}

	traceString, trace := captureBacktrace()
	return makeException(trace, code, msg, traceString)
}

// FromError creates an Exception from an error interface.
// If the error is already an Exception, it will be returned as is.
// This ensures that the error is always an Exception.
func FromError(err error) *Exception {
	var ex *Exception
	if errors.As(err, &ex) {
		return ex
	}
	traceString, trace := captureBacktrace()
	code := ""
	msg := err.Error()
	return makeException(trace, code, msg, traceString)
}

// Error method to satisfy the error interface.
func (e *Exception) Error() string {
	if e.previous != nil {
		return fmt.Sprintf(
			"Code: %v, Message: %s, Previous: %s",
			e.code,
			e.message,
			e.previous.Error(),
		)
	}
	return fmt.Sprintf("Code: %v, Message: %s", e.code, e.message)
}

// WithArgs captures arguments passed to the function where the Exception was raised. Useful for debugging purposes
func (e *Exception) WithArgs(args ...interface{}) *Exception {
	e.args = args
	return e
}

// GetCode returns the code of the exception.
func (e *Exception) GetCode() string {
	return e.code
}

// GetMessage returns the message of the exception.
func (e *Exception) GetMessage() string {
	return e.message
}

// GetPrevious returns the previous exception.
func (e *Exception) GetPrevious() *Exception {
	return e.previous
}

// GetTrace returns the backtrace of the exception.
func (e *Exception) GetTrace() *Trace {
	return e.trace
}

// GetTraceStr returns the backtrace of the exception as a string.
func (e *Exception) GetTraceStr() string {
	return e.traceStr
}

// GetFile returns the file name of the exception invocation.
func (e *Exception) GetFile() string {
	if e.origin == nil {
		return ""
	}
	return e.origin.File
}

// GetLine returns the line number of the exception invocation.
func (e *Exception) GetLine() uint {
	if e.origin == nil {
		return 0
	}
	return e.origin.Line
}

// GetPackage returns the package name of the exception invocation.
func (e *Exception) GetPackage() string {
	if e.origin == nil {
		return ""
	}
	return e.origin.Package
}

// GetFunction returns the function name of the exception invocation.
func (e *Exception) GetFunction() string {
	if e.origin == nil {
		return ""
	}
	return e.origin.Func
}

// GetStruct returns the struct name of the exception invocation.
func (e *Exception) GetStruct() string {
	if e.origin == nil {
		return ""
	}
	return e.origin.Struct
}

// GetField gets a custom field from the exception, as previously set in user-land.
func (e *Exception) GetField(key string) interface{} {
	return e.fields[key]
}

// JSON method to marshal the exception into a JSON string.
func (e *Exception) JSON() (string, error) {
	var toJSON func(ex *Exception) *exceptionJSON
	toJSON = func(ex *Exception) *exceptionJSON {
		if ex == nil {
			return nil
		}
		traceLines := ex.GetTrace().Lines
		return &exceptionJSON{
			Code:     ex.code,
			Message:  ex.message,
			Origin:   ex.origin,
			Trace:    traceLines,
			Previous: toJSON(ex.previous),
		}
	}

	exJSON := toJSON(e)
	bytes, err := json.Marshal(exJSON)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// SetPrevious sets a previous (chained) exception on the exception.
// Accepts an error interface, which can be a string or another Exception.
func (e *Exception) SetPrevious(previous error) *Exception {
	if previous != nil {
		var exception *Exception
		if errors.As(previous, &exception) {
			e.previous = exception
		} else {
			ex := &Exception{
				code:     "",
				message:  previous.Error(),
				previous: nil,
				traceStr: "",
				origin:   nil,
			}
			e.previous = ex
		}
	}
	return e
}

// SetField sets a custom field on the exception.
func (e *Exception) SetField(key string, value interface{}) *Exception {
	e.fields[key] = value
	return e
}

// SetFields sets multiple custom fields on the exception.
func (e *Exception) SetFields(fields Fields) *Exception {
	for key, value := range fields {
		e.fields[key] = value
	}
	return e
}

// Print prints a nicely formatted exception.
// For backtraceLevel, pass -1 to skip, pass 0 to print all levels, or pass a positive number to print up to that level.
// For chainLevel, pass -1 to skip, 0 to print all levels, or a positive number to print up to that level.
func (e *Exception) Print(backtraceLevel, chainLevel int) {
	isTTY := isatty.IsTerminal(os.Stdout.Fd())
	color := func(c string) string {
		if isTTY {
			return c
		}
		return ""
	}
	reset := func() string {
		if isTTY {
			return console.ColorReset
		}
		return ""
	}

	printException(e, backtraceLevel, chainLevel, 0, color, reset)
}

package e

import (
	"encoding/json"
	"errors"
	"fmt"
	"fx-service/pkg/console"
	"path/filepath"
	"runtime"
	"strings"
)

type exceptionJSON struct {
	Code     string         `json:"code"`
	Message  string         `json:"message"`
	Origin   *TraceLine     `json:"origin"`
	Trace    []TraceLine    `json:"trace"`
	Previous *exceptionJSON `json:"previous,omitempty"`
}

// trimFilePath trims the file path to include only the relevant parts.
func trimFilePath(fullPath string) string {
	index := strings.Index(fullPath, filePathTrimBefore)
	if index != -1 {
		start := index + len(filePathTrimBefore)
		return fullPath[start:]
	}
	return fullPath
}

// trimFunctionName splits the function string into package, struct, and function components.
func trimFunctionName(function string) (pkg string, structName string, fn string) {
	// Split the function string by "/"
	parts := strings.Split(function, "/")
	// The last part contains the package and function/method/lambda part
	lastPart := parts[len(parts)-1]

	// Extract just the package name from the last part before the first dot
	dotIndex := strings.Index(lastPart, ".")
	if dotIndex != -1 {
		pkg = lastPart[:dotIndex]
	} else {
		pkg = lastPart
	}

	// Check if there's a struct or not by looking for "(*StructName)"
	if strings.Contains(lastPart, "(*") {
		// Extract struct and method name
		start := strings.Index(lastPart, "(*") + 2
		end := strings.Index(lastPart, ")")
		structName = lastPart[start:end]
		fn = lastPart[end+2:]
	} else {
		// No struct, so extract the function name
		fn = lastPart[dotIndex+1:]
		structName = ""
	}

	return pkg, structName, fn
}

//func trimFunctionName(function string) string {
//	parts := strings.Split(function, "/")
//	return parts[len(parts)-1]
//}

// captureBacktrace captures the backtrace and returns the relevant details.
func captureBacktrace() (string, string, string, int, Trace) {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:n])

	var sb strings.Builder
	var file, function string
	var line int
	var traceLines []TraceLine
	for i := 0; ; i++ {
		frame, more := frames.Next()
		sb.WriteString(fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		pk, st, fn := trimFunctionName(frame.Function)
		trimmedFile := trimFilePath(frame.File)
		traceLines = append(traceLines, TraceLine{
			File:    trimmedFile,
			Line:    uint(frame.Line),
			Package: pk,
			Struct:  st,
			Func:    fn,
		})
		if i == 0 {
			function = fn
			file = trimmedFile
			line = frame.Line
		}
		if !more {
			break
		}
	}

	if file == "" {
		file = "[unknown file]"
	}
	if function == "" {
		function = "[unknown function]"
	}
	if line == 0 {
		line = 0 // Unknown line
	}

	return sb.String(), file, function, line, Trace{Lines: traceLines}
}

// makeException internal helper to create an exception instance.
func makeException(trace Trace, line int, code, msg, ts, file, function string) *Exception {
	tl := &TraceLine{
		File: file,
		Line: uint(line),

		Func: function,
	}

	return &Exception{
		code:     code,
		message:  msg,
		traceStr: ts,
		origin:   tl,
		trace:    &trace,
		fields:   make(map[string]interface{}),
	}
}

//func makeException(trace Trace, line int, code, msg, ts, file, function string) *Exception {
//	tl := &TraceLine{
//		File: file,
//		Line: uint(line),
//		Func: function,
//	}
//
//	return &Exception{
//		code:     code,
//		message:  msg,
//		traceStr: ts,
//		origin:   tl,
//		trace:    &trace,
//		fields:   make(map[string]interface{}),
//	}
//}

// printException private helper to recursively print exceptions.
func printException(ex *Exception, backtraceLevel, chainLevel, level int, color func(string) string, reset func() string) {
	fmt.Printf("%sException: '%s%s'\n", color(console.ColorRed), ex.message, reset())
	fmt.Printf("  • %sCode: %v%s\n", color(console.ColorCyan), ex.code, reset())
	fmt.Printf("  • %sFunction: %s%s\n", color(console.ColorGreen), ex.GetFunction(), reset())
	//fmt.Printf("  • %sIn file: %s, line:%d%s\n", color(console.ColorGreen), ex.GetFile(), ex.GetLine(), reset())

	// Extract the file name from the full file path
	fileName := filepath.Base(ex.GetFile())
	fullFilePath := strings.TrimPrefix(ex.GetFile(), "/")
	lineNumber := ex.GetLine()
	fmt.Printf("  • %sIn file: %s:%d [%s]%s\n", color(console.ColorGreen), fileName, lineNumber, fullFilePath, reset())

	// Check if we have custom fields present. If so, print a json string of them
	if len(ex.fields) > 0 {
		fieldsJSON, err := json.Marshal(ex.fields)
		if err == nil {
			fmt.Printf("  • %sFields: %s%s\n", color(console.ColorMagenta), string(fieldsJSON), reset())
		} else {
			fmt.Printf("[could not print fields] %v\n", err)
		}
	}

	// Print backtrace information if available
	// -1 to skip, 0 to print all, or a positive number to print up to that level
	if backtraceLevel >= 0 {
		trace := ex.GetTrace()
		if trace != nil {
			fmt.Printf("  • %sBacktrace:%s\n", color(console.ColorYellow), reset())
			end := trace.Count()
			if backtraceLevel > 0 && backtraceLevel < trace.Count() {
				end = backtraceLevel
			}
			for i := 0; i < end; i++ {
				traceLine := trace.Item(i)
				fmt.Printf("\t[%d]\t%s%s%s\n", i, color(console.ColorGrey), traceLine.String(), reset())
			}
		}
	}

	// Print previous exceptions if available
	// -1 to skip, 0 to print all, or a positive number to print up to that level
	if chainLevel >= 0 && ex.previous != nil {
		fmt.Printf("%sPrevious:%s\n", color(console.ColorMagenta), reset())
		if chainLevel > 0 && level >= chainLevel {
			return
		}
		var prevEx *Exception
		if errors.As(ex.previous, &prevEx) {
			printException(prevEx, backtraceLevel, chainLevel, level+1, color, reset)
		}
	}
}

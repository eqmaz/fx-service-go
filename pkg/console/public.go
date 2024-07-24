package console

// === Public Functions ------------------------------------------------------------------------------------------------

import "fmt"

// Suspend stops the console output until resumed. Useful for running tests
func Suspend() {
	suspendMu.Lock()
	suspended = true
	suspendMu.Unlock()
}

// Resume resumes the console output after being suspended
func Resume() {
	suspendMu.Lock()
	suspended = false
	suspendMu.Unlock()
}

// Out prints the formatted timestamp followed by the value
func Out(values ...interface{}) {
	send(valuesToStr(values))
}

// Outf prints the formatted timestamp followed by the formatted value
func Outf(format string, values ...interface{}) {
	send(fmt.Sprintf(format, values...))
}

// Success prints a success message with the timestamp
func Success(values ...interface{}) {
	send(colorize("✔ ", ColorGreen) + valuesToStr(values))
}

// Successf prints a formatted success message with the timestamp
func Successf(format string, value ...interface{}) {
	Success(fmt.Sprintf(format, value...))
}

// Warn prints a warning message with the timestamp
func Warn(values ...interface{}) {
	send(colorize("w ", ColorYellow) + valuesToStr(values))
}

// Warnf prints a formatted warning message with the timestamp
func Warnf(format string, value ...interface{}) {
	Warn(fmt.Sprintf(format, value...))
}

// Info prints an informational message with the timestamp
func Info(values ...interface{}) {
	send(colorize("ℹ ", ColorBlue) + valuesToStr(values))
}

// Infof prints a formatted informational message with the timestamp
func Infof(format string, value ...interface{}) {
	Info(fmt.Sprintf(format, value...))
}

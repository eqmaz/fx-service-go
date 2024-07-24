package console

import (
	"fmt"
	"github.com/mattn/go-isatty"
	"os"
	"sync"
	"time"
)

// === Private Variables -----------------------------------------------------------------------------------------------

// Output is only colored if the output destination is a terminal
var isATTY = isatty.IsTerminal(os.Stdout.Fd())

var suspended = false
var suspendMu sync.Mutex // Mutex to make Suspend and Resume thread-safe

// === Private Functions -----------------------------------------------------------------------------------------------
func send(str string) {
	suspendMu.Lock()
	defer suspendMu.Unlock()
	if !suspended {
		output := outTimestamp() + " " + str
		fmt.Println(output)
	}
}

// colorize wraps the given text with the provided color
func colorize(text, color string) string {
	if !isATTY {
		// If the output is not a terminal, return the text as is, without color
		return text
	}
	return fmt.Sprintf("%s%s%s", color, text, ColorReset)
}

// toStr converts a generic value to a string using fmt.Sprintf
func toStr(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// valuesToStr converts a slice of values to a string
func valuesToStr(values []interface{}) string {
	output := ""
	for i, v := range values {
		if i > 0 {
			output += " "
		}
		output += toStr(v)
	}
	return output
}

// outTimestamp returns a string representation of the current time in the format "[%Y-%m-%d %H:%M:%S]"
func outTimestamp() string {
	currentTime := time.Now()
	timeStr := fmt.Sprintf("[%s]", currentTime.Format("2006-01-02 15:04:05.000"))
	return colorize(timeStr, ColorGrey)
}

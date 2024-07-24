package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

// Dump prints data on the console in a pretty format
func Dump(data interface{}) {
	// Get the type of the data
	dataType := reflect.TypeOf(data)
	fmt.Printf("[%s]\n", dataType)

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return
	}
	fmt.Println(string(jsonData))
}

// Dd prints data on the console in a pretty format and exits the program
func Dd(data interface{}) {
	Dump(data)
	os.Exit(1)
}

// CaptureOutput captures the output of a function that writes to stdout
func CaptureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	outC := make(chan string)
	// Copy output to a buffer in a separate goroutine
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		outC <- buf.String()
	}()

	f()

	// Restore stdout and close the pipe
	_ = w.Close()
	os.Stdout = stdout
	return <-outC
}

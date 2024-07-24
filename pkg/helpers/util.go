// Package helpers Helpers
package helpers

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
)

// ExeDir returns the directory of the executable
func ExeDir() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not get executable path: %w", err)
	}
	return filepath.Dir(exePath), nil
}

// IsExeInCwd checks if the executable's directory is the current working directory.
func IsExeInCwd() (bool, error) {
	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return false, err
	}

	// Get the executable's path
	exePath, err := os.Executable()
	if err != nil {
		return false, err
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exePath)

	// Compare the current working directory with the executable's directory
	return filepath.Clean(cwd) == filepath.Clean(exeDir), nil
}

// IsReadableFile checks if a file exists and is readable
func IsReadableFile(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		return false
	}
	err = file.Close()
	if err != nil {
		return false
	}
	return true
}

// GetType returns the type of the variable as a string
func GetType(variable interface{}) string {
	return reflect.TypeOf(variable).String()
}

// Round rounds a number to a specified number of decimal places
func Round(num float64, precision int) float64 {
	factor := math.Pow(10, float64(precision))
	return math.Round(num*factor) / factor
}

// NumberFormat formats a number into a string, with decimal and thousands separators
func NumberFormat(number float64, decimals int, decPoint, thousandsSep string) string {
	format := fmt.Sprintf("%%.%df", decimals)
	formatted := fmt.Sprintf(format, number)

	// Split the number into integer and decimal parts
	parts := []byte(formatted)
	var intPart, decPart []byte
	pointFound := false
	for _, c := range parts {
		if c == '.' {
			pointFound = true
			continue
		}
		if pointFound {
			decPart = append(decPart, c)
		} else {
			intPart = append(intPart, c)
		}
	}

	// Insert thousands separator
	var intPartWithSep []byte
	for i, c := range intPart {
		if i > 0 && (len(intPart)-i)%3 == 0 {
			intPartWithSep = append(intPartWithSep, []byte(thousandsSep)...)
		}
		intPartWithSep = append(intPartWithSep, c)
	}

	// Join integer and decimal parts with the decimal point
	result := string(intPartWithSep)
	if decimals > 0 {
		result += decPoint + string(decPart)
	}

	return result
}

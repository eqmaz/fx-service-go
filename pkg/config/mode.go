package config

import (
	"encoding/json"
	"fmt"
)

// Mode type - this is the strategy for choosing the next provider API
type Mode int

// Define constants for each mode using iota
const (
	Race      Mode = iota // All healthy providers are called at the same time, the first one to respond wins
	Robin                 // Healthy Providers are called in a round-robin fashion
	First                 // The first healthy provider is called
	Random                // A random healthy provider is chosen for each request
	Priority              // Healthy providers are called in order of priority
	Aggregate             // All healthy providers are called, and the results are aggregated (averaged)
)

// ModeNameList returns a simple list of mode names
func ModeNameList() []string {
	// Construct the error message dynamically
	modes := []Mode{Race, Robin, First, Random, Priority, Aggregate}
	result := make([]string, len(modes))
	for i, mode := range modes {
		result[i] = mode.String()
	}
	return result
}

// String method to get the string representation of the mode
func (m *Mode) String() string {
	return [...]string{"race", "robin", "first", "random", "priority", "aggregate"}[*m]
}

// UnmarshalJSON method to convert JSON string to Mode type
func (m *Mode) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case "race":
		*m = Race
	case "robin":
		*m = Robin
	case "first":
		*m = First
	case "random":
		*m = Random
	case "priority":
		*m = Priority
	case "aggregate":
		*m = Aggregate
	default:
		return fmt.Errorf("unsupported mode value '%s'. Use one of: %s", s, ModeNameList())
	}
	return nil
}

package helpers

import "math/rand"

// SliceContains returns true if slice contains the givens string
func SliceContains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

// RemoveSliceElement safely removes an element from a slice, if it's in bounds
func RemoveSliceElement(slice []string, index int) []string {
	if index < 0 || index >= len(slice) {
		return slice
	}
	return append(slice[:index], slice[index+1:]...)
}

// GetRandomSliceIndex returns a random index for a slice
func GetRandomSliceIndex(slice []string) int {
	return rand.Intn(len(slice))
}

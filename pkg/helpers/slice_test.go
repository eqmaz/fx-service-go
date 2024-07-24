package helpers

import (
	"testing"
)

// TestSliceContains tests the SliceContains function
func TestSliceContains(t *testing.T) {
	tests := []struct {
		slice    []string
		str      string
		contains bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "b", "c"}, "d", false},
		{[]string{}, "a", false},
	}

	for _, test := range tests {
		result := SliceContains(test.slice, test.str)
		if result != test.contains {
			t.Errorf("SliceContains(%v, %s) = %v; want %v", test.slice, test.str, result, test.contains)
		}
	}
}

// TestRemoveSliceElement tests the RemoveSliceElement function
func TestRemoveSliceElement(t *testing.T) {
	tests := []struct {
		slice    []string
		index    int
		expected []string
	}{
		{[]string{"a", "b", "c"}, 1, []string{"a", "c"}},
		{[]string{"a", "b", "c"}, 3, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, -1, []string{"a", "b", "c"}},
		{[]string{"a", "b", "c"}, 0, []string{"b", "c"}},
		{[]string{"a"}, 0, []string{}},
	}

	for _, test := range tests {
		result := RemoveSliceElement(test.slice, test.index)
		if !equalSlices(result, test.expected) {
			t.Errorf("RemoveSliceElement(%v, %d) = %v; want %v", test.slice, test.index, result, test.expected)
		}
	}
}

// TestGetRandomSliceIndex tests the GetRandomSliceIndex function
func TestGetRandomSliceIndex(t *testing.T) {
	slice := []string{"a", "b", "c"}
	index := GetRandomSliceIndex(slice)
	if index < 0 || index >= len(slice) {
		t.Errorf("GetRandomSliceIndex(%v) = %d; want index in range [0, %d)", slice, index, len(slice))
	}

	// Testing empty slice edge case
	var emptySlice []string
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("GetRandomSliceIndex did not panic on empty slice")
		}
	}()
	GetRandomSliceIndex(emptySlice)
}

// Helper function to compare two slices for equality
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

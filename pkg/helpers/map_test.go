package helpers

import (
	"reflect"
	"sort"
	"testing"
)

// TestGetMapKeys tests the GetMapKeys function
func TestGetMapKeys(t *testing.T) {
	tests := []struct {
		name string
		m    interface{}
		want interface{}
	}{
		{"empty string map", map[string]int{}, []string{}},
		{"single key string map", map[string]int{"a": 1}, []string{"a"}},
		{"multiple keys string map", map[string]int{"a": 1, "b": 2, "c": 3}, []string{"a", "b", "c"}},
		{"empty int map", map[int]string{}, []int{}},
		{"single key int map", map[int]string{1: "a"}, []int{1}},
		{"multiple keys int map", map[int]string{1: "a", 2: "b"}, []int{1, 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch m := tt.m.(type) {
			case map[string]int:
				got := GetMapKeys(m)
				want := tt.want.([]string)
				sort.Strings(got)
				sort.Strings(want)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("GetMapKeys() = %v, want %v", got, want)
				}
			case map[int]string:
				got := GetMapKeys(m)
				want := tt.want.([]int)
				sort.Ints(got)
				sort.Ints(want)
				if !reflect.DeepEqual(got, want) {
					t.Errorf("GetMapKeys() = %v, want %v", got, want)
				}
			default:
				t.Fatalf("unsupported map type %T", tt.m)
			}
		})
	}
}

// TestCloneMapShallow tests the CloneMapShallow function
func TestCloneMapShallow(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{"nil input", nil, nil},
		{"empty map", map[string]int{}, map[string]int{}},
		{"single key", map[string]int{"a": 1}, map[string]int{"a": 1}},
		{"multiple keys", map[string]int{"a": 1, "b": 2}, map[string]int{"a": 1, "b": 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CloneMapShallow(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneMapShallow() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCloneMapDeep tests the CloneMapDeep function
func TestCloneMapDeep(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{"nil input", nil, nil},
		{"empty map", map[string]int{}, map[string]int{}},
		{"simple map", map[string]int{"a": 1}, map[string]int{"a": 1}},
		{"nested map", map[string]map[string]int{"a": {"b": 2}}, map[string]map[string]int{"a": {"b": 2}}},
		{"map with slice", map[string][]int{"a": {1, 2}}, map[string][]int{"a": {1, 2}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CloneMapDeep(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CloneMapDeep() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Additional test to ensure deep copy is truly deep
func TestCloneMapDeepNested(t *testing.T) {
	input := map[string]map[string]int{"a": {"b": 2}}
	got := CloneMapDeep(input).(map[string]map[string]int)
	got["a"]["b"] = 3

	if input["a"]["b"] == got["a"]["b"] {
		t.Errorf("CloneMapDeep() did not create a deep copy")
	}
}

package e

import (
	"reflect"
	"testing"
)

func TestFields_With(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value interface{}
		start Fields
		want  Fields
	}{
		{
			name:  "add string value",
			key:   "key1",
			value: "value1",
			start: Fields{},
			want:  Fields{"key1": "value1"},
		},
		{
			name:  "add int value",
			key:   "key2",
			value: 42,
			start: Fields{},
			want:  Fields{"key2": 42},
		},
		{
			name:  "add to existing fields",
			key:   "key3",
			value: true,
			start: Fields{"key1": "value1"},
			want:  Fields{"key1": "value1", "key3": true},
		},
		{
			name:  "overwrite existing key",
			key:   "key1",
			value: "new_value",
			start: Fields{"key1": "old_value"},
			want:  Fields{"key1": "new_value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := tt.start
			if got := f.With(tt.key, tt.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fields.With() = %v, want %v", got, tt.want)
			}
		})
	}
}

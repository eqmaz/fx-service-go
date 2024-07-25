package e

import (
	"testing"
)

func TestTraceLine_String(t *testing.T) {
	tests := []struct {
		name      string
		traceLine TraceLine
		expected  string
	}{
		{"basic case", TraceLine{File: "file.go", Line: 10, Package: "mkPkg", Struct: "foo", Func: "FuncName"}, "file.go:10 mkPkg.foo.FuncName"},
		{"empty values", TraceLine{File: "", Line: 0, Package: "", Struct: "", Func: ""}, ":0 ."},
		{"no struct", TraceLine{File: "file.go", Line: 18, Package: "foo", Func: "Bar"}, "file.go:18 foo.Bar"},
		{
			"long function name",
			TraceLine{File: "file.go", Line: 123456, Package: "mkPkg", Struct: "foo", Func: "AReallyLongFunctionName"},
			"file.go:123456 mkPkg.foo.AReallyLongFunctionName",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.traceLine.String()
			if result != tt.expected {
				t.Errorf("\nexpected: %q, \ngot: %q", tt.expected, result)
			}
		})
	}
}

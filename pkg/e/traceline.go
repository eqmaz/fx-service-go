package e

import (
	"fmt"
)

// TraceLine represents a single line (or step) in the backtrace.
// A trace has multiple TraceLines.
type TraceLine struct {
	File    string
	Line    uint
	Package string
	Struct  string
	Func    string
}

// String returns a string representation of the TraceLine.
func (t TraceLine) String() string {
	st := ""
	if t.Struct != "" {
		st = t.Struct + "."
	}
	return fmt.Sprintf("%s:%d %s.%s%s", t.File, t.Line, t.Package, st, t.Func)
}

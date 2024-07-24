package e

// Trace struct to represent a collection of TraceLines.
type Trace struct {
	Lines []TraceLine
}

// Top returns the first TraceLine in the Trace.
func (t *Trace) Top() *TraceLine {
	if len(t.Lines) == 0 {
		return nil
	}
	return &t.Lines[0]
}

// Last gets the last TraceLine in the Trace.
func (t *Trace) Last() *TraceLine {
	if len(t.Lines) == 0 {
		return nil
	}
	return &t.Lines[len(t.Lines)-1]
}

// Item gets the TraceLine at the given index.
func (t *Trace) Item(n int) *TraceLine {
	if n < 0 || n >= len(t.Lines) {
		return nil
	}
	return &t.Lines[n]
}

// Count the number of TraceLines in the Trace.
func (t *Trace) Count() int {
	return len(t.Lines)
}

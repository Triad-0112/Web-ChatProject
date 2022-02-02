package tracer

import (
	"fmt"
	"io"
)

//tracer method
type Tracer interface {
	Trace(...interface{})
}

//traced struct
type tracer struct {
	out io.Writer
}

//unnecessary tracer struct
type nilTracer struct{}

//unnecessary tracer method
func (t *nilTracer) Trace(a ...interface{}) {}

// Function to Avoid noisy unnecessary trace
func Off() Tracer {
	return &nilTracer{}
}

//Function to Print out every activity
func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

//Create allocated tracer
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

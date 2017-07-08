package finger

import (
	"io"
)

// Handler responds to a Finger query.
//
// ServerFinger should respond to a finger query and write the response
// to the given io.Writer.
type Handler interface {
	ServeFinger(io.Writer, *Query)
}

// HandlerFunc allows use of ordinary functions as finger handlers.
type HandlerFunc func(io.Writer, *Query)

func (f HandlerFunc) ServeFinger(w io.Writer, q *Query) {
	f(w, q)
}

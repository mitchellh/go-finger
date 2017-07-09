package finger

import (
	"context"
	"io"
)

// Handler responds to a Finger query.
//
// ServerFinger should respond to a finger query and write the response
// to the given io.Writer.
type Handler interface {
	ServeFinger(context.Context, io.Writer, *Query)
}

// HandlerFunc allows use of ordinary functions as finger handlers.
type HandlerFunc func(context.Context, io.Writer, *Query)

func (f HandlerFunc) ServeFinger(ctx context.Context, w io.Writer, q *Query) {
	f(ctx, w, q)
}

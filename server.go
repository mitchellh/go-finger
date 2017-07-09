package finger

import (
	"bufio"
	"context"
	"io"
	"net"
	"time"
)

// Server defines parameters for running a Finger server. The zero value
// of a Server is a valid configuration, though every request will be
// closed immediately since no handler is set.
type Server struct {
	Addr    string
	Handler Handler

	// ReadTimeout is the maximum duration before timing out reads of the
	// response. This sets a deadline on the connection and isn't a handler
	// timeout.
	ReadTimeout time.Duration

	// WriteTimeout is the maximum duration before timing out writes of the
	// response. This sets a deadline on the connection and isn't a handler
	// timeout.
	WriteTimeout time.Duration

	// MaxQueryBytes is the maximum amount of bytes that will be read from
	// the connection to determine the query.
	MaxQueryBytes int
}

// Serve listens on the finger port and serves connections with the given
// handler. This sets reasonable read/write timeouts. Please see the source
// for the exact timeouts which should be generously high.
func Serve(h Handler) error {
	s := &Server{
		Handler:       h,
		ReadTimeout:   5 * time.Minute,
		WriteTimeout:  5 * time.Minute,
		MaxQueryBytes: 4096,
	}

	return s.ListenAndServe()
}

// ListenAndServe listens on the TCP network address s.Addr and then
// calls Serve to handle incoming connections. If s.Addr is blank,
// ":finger" is used.
func (s *Server) ListenAndServe() error {
	addr := s.Addr
	if addr == "" {
		addr = ":finger"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return s.Serve(ln)
}

// Serve accepts incoming connections on the listener l, creating a new
// service goroutine for each.
//
// The listener is closed when this function returns.
func (s *Server) Serve(l net.Listener) error {
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}

		s.setupConn(c)
		go s.ServeConn(context.Background(), c)
	}
}

// ServeConn serves a single connection, blocking until the request is complete.
func (s *Server) ServeConn(ctx context.Context, conn io.ReadWriteCloser) error {
	// Close the connection once we're done in any case
	defer conn.Close()

	// If we have a maximum amount to read, then setup a limit
	var r io.Reader = conn
	if s.MaxQueryBytes > 0 {
		r = io.LimitReader(conn, int64(s.MaxQueryBytes))
	}

	// Read the query line
	buf := bufio.NewReader(r)
	line, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	// Parse the query
	query, err := ParseQuery(line)
	if err != nil {
		return err
	}

	// If we have a net conn then setup the remote addr
	if nc, ok := conn.(net.Conn); ok {
		query.RemoteAddr = nc.RemoteAddr()
	}

	if s.Handler != nil {
		s.Handler.ServeFinger(ctx, conn, query)
	}

	return nil
}

func (s *Server) setupConn(c net.Conn) {
	t0 := time.Now()
	if s.ReadTimeout > 0 {
		c.SetReadDeadline(t0.Add(s.ReadTimeout))
	}

	if s.WriteTimeout > 0 {
		c.SetWriteDeadline(t0.Add(s.WriteTimeout))
	}
}

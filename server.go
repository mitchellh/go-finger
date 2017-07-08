package finger

import (
	"bufio"
	"context"
	"io"
	"net"
)

type Server struct {
	Addr    string
	Handler Handler
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

		go s.ServeConn(context.Background(), c)
	}
}

// ServeConn serves a single connection, blocking until the request is complete.
func (s *Server) ServeConn(ctx context.Context, conn io.ReadWriteCloser) error {
	// Close the connection once we're done in any case
	defer conn.Close()

	// Read the query line
	buf := bufio.NewReader(conn)
	line, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	// Parse the query
	query, err := ParseQuery(line)
	if err != nil {
		return err
	}

	s.Handler.ServeFinger(conn, query)
	return nil
}

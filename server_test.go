package finger

import (
	"bytes"
	"context"
	"io"
	"net"
	"os/exec"
	"testing"
)

func TestServerServeConn(t *testing.T) {
	ln := testFingerLn(t)
	defer ln.Close()

	// Setup the handler
	handler := func(ctx context.Context, w io.Writer, q *Query) {
		w.Write([]byte("received!"))
	}

	// Serve
	s := &Server{Handler: HandlerFunc(handler)}
	go s.Serve(ln)

	// Use the finger command to run it
	actual := testFinger(t, "foo@127.0.0.1")
	if !bytes.Contains(actual, []byte("received!")) {
		t.Fatalf("bad: %s", actual)
	}
}

func TestServerServeConn_maxQueryBytes(t *testing.T) {
	ln := testFingerLn(t)
	defer ln.Close()

	// Setup the handler
	handler := func(ctx context.Context, w io.Writer, q *Query) {
		panic("don't call me!")
	}

	// Serve
	s := &Server{Handler: HandlerFunc(handler)}
	s.MaxQueryBytes = 4
	go s.Serve(ln)

	// Use the finger command to run it
	testFinger(t, "toomanybytes@127.0.0.1")
}

func testFingerLn(t *testing.T) net.Listener {
	ln, err := net.Listen("tcp", ":finger")
	if err != nil {
		t.Skipf("can't listen on finger port: %s", err)
		t.SkipNow()
	}

	return ln
}

func testFinger(t *testing.T, query string) []byte {
	if _, err := exec.LookPath("finger"); err != nil {
		t.Skipf("finger not found")
		t.SkipNow()
	}

	actual, err := exec.Command("finger", query).Output()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return actual
}

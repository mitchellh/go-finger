package finger

import (
	"bytes"
	"io"
	"net"
	"os"
	"os/exec"
	"testing"
)

func TestMain(m *testing.M) {
	// Verify we can listen on the finger port
	ln, err := net.Listen("tcp", ":finger")
	if err != nil {
		println("Must run as root to be able to listen on port 79 for tests", err.Error())
		os.Exit(1)
	}
	ln.Close()

	os.Exit(m.Run())
}

func TestServerServeConn(t *testing.T) {
	// Create a listener
	ln, err := net.Listen("tcp", ":finger")
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	defer ln.Close()

	// Setup the handler
	handler := func(w io.Writer, q *Query) {
		w.Write([]byte("received!"))
	}

	// Serve
	s := &Server{Handler: HandlerFunc(handler)}
	go s.Serve(ln)

	// Use the finger command to run it
	actual, err := exec.Command("finger", "foo@127.0.0.1").Output()
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if !bytes.Contains(actual, []byte("received!")) {
		t.Fatalf("bad: %s", actual)
	}
}

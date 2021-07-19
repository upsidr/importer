package stdout

import (
	"io"
	"os"
	"testing"
)

type Stdout struct {
	origStdout *os.File

	fakeStdin  *os.File
	fakeStdout *os.File

	isClosed bool
}

func New(t testing.TB) *Stdout {
	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	return &Stdout{
		origStdout: origStdout,
		fakeStdin:  r,
		fakeStdout: w,
	}
}

// Read reads from Stdout. After reading once, Stdout setup will be closed and
// thus cannot be reused.
func (s *Stdout) Read(t testing.TB) []byte {
	if s.isClosed {
		t.Fatalf("stdout is already closed")
	}
	s.fakeStdout.Close()
	os.Stdout = s.origStdout

	stdoutResult, err := io.ReadAll(s.fakeStdin)
	if err != nil {
		t.Fatal(err)
	}
	s.fakeStdin.Close()

	s.isClosed = true
	return stdoutResult
}

// Close closes all open files for Stdout testing, and put back the original
// stdout in place.
func (s *Stdout) Close() {
	s.fakeStdout.Close()
	s.fakeStdin.Close()
	os.Stdout = s.origStdout
	s.isClosed = true
}

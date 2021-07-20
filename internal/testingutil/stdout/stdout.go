package stdout

import (
	"io"
	"os"
	"testing"
)

// Stdout holds fake Stdout, which can be used to replace normal Stdout for
// testing. Use New() function to create this struct.
type Stdout struct {
	origStdout   *os.File
	stdoutWriter *os.File
	stdoutReader *os.File

	isClosed bool
}

// New sets up os.File for os.Stdout, and ensures any stdout is captured.
// Captured output can be read by Stdout.Read(testing.TB) function.
func New(t testing.TB) *Stdout {
	t.Helper()

	origStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	os.Stdout = w

	return &Stdout{
		origStdout:   origStdout,
		stdoutWriter: w,
		stdoutReader: r,
	}
}

// ReadAllAndClose reads from Stdout. After reading once, Stdout setup will be
// closed and thus cannot be reused.
func (s *Stdout) ReadAllAndClose(t testing.TB) []byte {
	t.Helper()

	if s.isClosed {
		t.Fatalf("stdout is already closed")
	}

	s.stdoutWriter.Close()
	os.Stdout = s.origStdout

	stdoutResult, err := io.ReadAll(s.stdoutReader)
	if err != nil {
		t.Fatal(err)
	}
	s.stdoutReader.Close()

	s.isClosed = true
	return stdoutResult
}

// Close closes all open files for Stdout testing, and put back the original
// stdout in place.
func (s *Stdout) Close() {
	s.stdoutWriter.Close()
	s.stdoutReader.Close()
	os.Stdout = s.origStdout

	s.isClosed = true
}

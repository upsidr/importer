package golden

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

// CopyTemp reads in the file from given path, and creates a copy file
// "in the same" directory. The name of temporarily copied file is a random
// string, with extension string as suffix. If you are providing an extension,
// make sure to include '.' as well, e.g. '.md', '.yaml'.
//
// The first return value is the file path, and the second value is a clean up
// function to remove the file.
func CopyTemp(t testing.TB, path string) (string, func()) {
	in, err := os.Open(path)
	if err != nil {
		t.Fatalf("error with reading file '%s', %v", path, err)
	}
	defer in.Close()

	dir := filepath.Dir(path)
	ext := filepath.Ext(path)

	out, err := os.CreateTemp(dir, "copy_*"+ext)
	if err != nil {
		t.Fatalf("error with creating a temp file, %v", err)
	}
	defer out.Close()
	dstFileName := out.Name()

	_, err = io.Copy(out, in)
	if err != nil {
		t.Fatalf("copy from '%s' to '%s' failed, %v", path, dstFileName, err)
	}

	return dstFileName, func() { os.Remove(dstFileName) }
}

// File reads content from provided path, and returns byte representation of
// the file while handling error with testing.TB.
func File(t testing.TB, path string) []byte {
	t.Helper()
	x, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error reading file '%s', %v", path, err)
	}
	return x
}

// FileAsString reads content from provided path, and returns string
// representation of the file while handling error with testing.TB.
func FileAsString(t testing.TB, path string) string {
	t.Helper()
	b := File(t, path)
	return string(b)
}

// UpdateFile updates the golden file at the provided path with content. This
// fails if there is no file at the given path.
func UpdateFile(t testing.TB, path string, content []byte) {
	t.Helper()
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		t.Fatalf("failed to write to golden file '%s', %v", path, err)
	}
}

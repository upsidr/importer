package file

import (
	"os"
	"strings"
)

type replaceMode struct {
	isDryRun bool
	isForce  bool
}
type ReplaceOption func(*replaceMode) error

func WithDryRun() ReplaceOption {
	return func(m *replaceMode) error {
		m.isDryRun = true
		return nil
	}
}

func WithForce() ReplaceOption {
	return func(m *replaceMode) error {
		m.isForce = true
		return nil
	}
}

// ReplaceWithAfter replaces the original file content with the processed
// content. This is done by creating a temp file first, and replacing it.
//
// TODO: Ensure file mode is kept, or clarify in the comment.
func (f *File) ReplaceWithAfter(options ...ReplaceOption) error {
	mode := &replaceMode{}
	for _, opt := range options {
		opt(mode)
	}

	file, err := os.CreateTemp("/tmp/", "importer_replace_*")
	if err != nil {
		return err // TODO: test coverage
	}
	defer file.Close()

	_, err = file.Write(f.ContentAfter)
	if err != nil {
		return err // TODO: test coverage
	}

	// Exit early if dry-run
	if mode.isDryRun {
		return nil
	}

	err = os.Rename(file.Name(), f.FileName)
	if err != nil {
		return err // TODO: test coverage
	}

	return nil
}

// ReplaceWithAfter replaces the original file content with the processed
// content. This is done by creating a temp file first, and replacing it.
//
// TODO: Ensure file mode is kept, or clarify in the comment.
func (f *File) ReplaceWithPurged(options ...ReplaceOption) error {
	mode := &replaceMode{}
	for _, opt := range options {
		opt(mode)
	}

	file, err := os.CreateTemp("/tmp/", "importer_replace_*")
	if err != nil {
		return err // TODO: test coverage
	}
	defer file.Close()

	data := strings.Join(f.ContentPurged, "\n")
	data = data + "\n" // Make sure to add new line at the end of the file
	_, err = file.WriteString(data)
	if err != nil {
		return err // TODO: test coverage
	}

	// Exit early if dry-run
	if mode.isDryRun {
		return nil
	}

	err = os.Rename(file.Name(), f.FileName)
	if err != nil {
		return err // TODO: test coverage
	}

	return nil
}

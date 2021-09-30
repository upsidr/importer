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
// content.
func (f *File) ReplaceWithAfter(options ...ReplaceOption) error {
	mode := &replaceMode{}
	for _, opt := range options {
		opt(mode)
	}

	// If SkipUpdate flag is on, do not update the file.
	if f.SkipUpdate {
		// Perhaps it's more friendly to log that the file update is being
		// skipped. If that's truly the case, handle it here.
		mode.isDryRun = true
	}

	return replace(f.FileName, f.ContentAfter, mode)
}

// ReplaceWithAfter replaces the original file content with the processed
// content.
func (f *File) ReplaceWithPurged(options ...ReplaceOption) error {
	mode := &replaceMode{}
	for _, opt := range options {
		opt(mode)
	}

	data := strings.Join(f.ContentPurged, "\n")
	data = data + "\n" // Make sure to add new line at the end of the file

	return replace(f.FileName, []byte(data), mode)
}

func replace(fileName string, content []byte, mode *replaceMode) error {
	// Exit early if dry-run
	if mode.isDryRun {
		return nil
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(content)
	return err
}

package file

import (
	"os"
	"strings"
)

// ReplaceWithAfter replaces the original file content with the processed
// content. This is done by creating a temp file first, and replacing it.
//
// TODO: Ensure file mode is kept, or clarify in the comment.
func (f *File) ReplaceWithAfter() error {
	file, err := os.CreateTemp("/tmp/", "importer_replace_*")
	if err != nil {
		return err // TODO: test coverage
	}
	defer file.Close()

	_, err = file.Write(f.ContentAfter)
	if err != nil {
		return err // TODO: test coverage
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
func (f *File) ReplaceWithPurged() error {
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

	err = os.Rename(file.Name(), f.FileName)
	if err != nil {
		return err // TODO: test coverage
	}

	return nil
}

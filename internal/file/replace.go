package file

import (
	"bytes"
	"os"
)

// ReplaceWithAfter replaces the original file content with the processed
// content. This is done by creating a temp file first, and replacing it.
//
// TODO: Ensure file mode is kept, or clarify in the comment.
func (f *File) ReplaceWithAfter() error {
	file, err := os.CreateTemp("/tmp/", "importer_replace_*")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(f.ContentAfter)
	if err != nil {
		return err
	}

	err = os.Rename(file.Name(), f.FileName)
	if err != nil {
		return err
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
		return err
	}
	defer file.Close()

	data := bytes.Join(f.ContentPurged, []byte("\n"))
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	err = os.Rename(file.Name(), f.FileName)
	if err != nil {
		return err
	}

	return nil
}

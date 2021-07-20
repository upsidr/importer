package file

import "os"

// WriteAfterTo takes the processed content and copies into provided filepath.
//
// TODO: Ensure file mode is kept, or clarify in the comment.
func (f *File) WriteAfterTo(filepath string) error {
	file, err := os.CreateTemp("/tmp/", "importer_replace_*")
	if err != nil {
		return err // TODO: test coverage
	}
	defer file.Close()

	_, err = file.Write(f.ContentAfter)
	if err != nil {
		return err // TODO: test coverage
	}

	err = os.Rename(file.Name(), filepath)
	if err != nil {
		return err // TODO: test coverage
	}

	return nil
}

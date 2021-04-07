package file

import "os"

// ReplaceWithImporter replaces the original file content with the processed
// content. This is done by creating a temp file first, and replacing it.
//
// TODO: Ensure file mode is kept, or clarify in the comment.
func (f *File) ReplaceWithImporter() error {
	file, err := os.CreateTemp("/tmp/", "importer_replace_*")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(f.contentAfter)
	if err != nil {
		return err
	}

	err = os.Rename(file.Name(), f.FileName)
	if err != nil {
		return err
	}

	return nil
}

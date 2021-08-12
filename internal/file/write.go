package file

import "os"

// WriteAfterTo writes the processed content to the provided filepath.
func (f *File) WriteAfterTo(filepath string) error {
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(f.ContentAfter)
	return err
}

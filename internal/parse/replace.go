package parse

import "os"

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

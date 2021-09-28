package file

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteAfterTo writes the processed content to the provided filepath.
func (f *File) WriteAfterTo(targetFilePath string) error {
	file, err := os.OpenFile(targetFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	content := []byte{}
	content = append(content, f.prepareGeneratedHeader(targetFilePath)...)
	content = append(content, f.ContentAfter...)

	_, err = file.Write(content)
	return err
}

func (f *File) prepareGeneratedHeader(targetFilePath string) []byte {
	fileType := filepath.Ext(f.FileName)

	comment := `# == %s ==`
	switch fileType {
	case ".md":
		comment = `<!-- == %s == -->`
	case ".yaml", ".yml":
		comment = `# == %s ==`
	default:
	}

	// Compare directories of each file, rather than file itself, so that we
	// don't end up with an extra "../".
	relDir, err := filepath.Rel(filepath.Dir(targetFilePath), filepath.Dir(f.FileName))
	if err != nil {
		// This shouldn't happen, but when it does, simply ignore and return an
		// empty slice.
		return nil
	}

	baseFile := fmt.Sprintf("%s/%s", relDir, filepath.Base(f.FileName))
	note := fmt.Sprintf(`improter-generated-from: %s`, baseFile)

	x := fmt.Sprintf(comment+"\n", note)

	return []byte(x)
}

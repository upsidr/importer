package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func (f *File) UpdateWithAnnotations() error {
	result := []byte{}
	br := byte('\n')
	for line, data := range f.ContentPurged {
		result = append(result, data...)
		result = append(result, br)

		if a, found := f.Annotations[line+1]; found {
			// Make sure the files are read based on the relative path
			dir := filepath.Dir(f.FileName)
			targetPath := dir + "/" + a.TargetPath
			f, err := os.Open(targetPath)
			if err != nil {
				fmt.Printf("warning: could not open file '%s', skipping\n", targetPath)
				continue
			}
			currentLine := 0
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				currentLine++
				for _, l := range a.TargetLines {
					if currentLine == l {
						result = append(result, scanner.Bytes()...)
						result = append(result, br)
					}
				}
			}
		}
	}
	f.ContentAfter = result
	return nil
}

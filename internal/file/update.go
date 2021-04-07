package file

import (
	"bufio"
	"fmt"
	"os"
)

func (f *File) UpdateWithAnnotations() error {
	result := []byte{}
	br := byte('\n')
	for line, data := range f.ContentPurged {
		result = append(result, data...)
		result = append(result, br)

		if a, found := f.Annotations[line+1]; found {
			// TODO: file should be relative to the original file rather than current dir
			f, err := os.Open(a.TargetPath)
			if err != nil {
				fmt.Printf("warning: could not open file '%s', skipping\n", a.TargetPath)
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

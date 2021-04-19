package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// ProcessAnnotations reads annotations and generates ContentAfter.
func (f *File) ProcessAnnotations() error {
	result := []byte{}
	br := byte('\n')
	for line, data := range f.ContentPurged {
		result = append(result, data...)
		result = append(result, br)

		if a, found := f.Annotations[line+1]; found {
			// Make sure the files are read based on the relative path
			dir := filepath.Dir(f.FileName)
			targetPath := dir + "/" + a.TargetPath
			file, err := os.Open(targetPath)
			if err != nil {
				fmt.Printf("warning: could not open file '%s', skipping\n", targetPath)
				continue
			}
			defer file.Close() // TODO: Move logic within this for loop to separate func, so that defer runs as early as possible

			// Handle marker imports

			// Handle line number imports
			currentLine := 0
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				currentLine++
				if currentLine >= a.TargetLineFrom && currentLine <= a.TargetLineTo {
					result = append(result, scanner.Bytes()...)
					result = append(result, br)
					continue
				}
				for _, l := range a.TargetLines {
					if currentLine == l {
						result = append(result, scanner.Bytes()...)
						result = append(result, br)
						continue
					}
				}
			}
		}
	}
	f.ContentAfter = result
	return nil
}

func (f *File) processSingleAnnotation(result []byte, annotation *Annotation) ([]byte, error) {
	br := byte('\n')
	// Make sure the files are read based on the relative path
	dir := filepath.Dir(f.FileName)
	targetPath := dir + "/" + annotation.TargetPath
	file, err := os.Open(targetPath)
	if err != nil {
		// Purposely returning the byte slice as it contains data that were
		// populated prior to hitting this func
		return result, fmt.Errorf("could not open file '%s', skipping, %v", targetPath, err)
	}
	defer file.Close()

	currentLine := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++
		if currentLine >= annotation.TargetLineFrom &&
			currentLine <= annotation.TargetLineTo {
			result = append(result, scanner.Bytes()...)
			result = append(result, br)
			continue
		}
		for _, l := range annotation.TargetLines {
			if currentLine == l {
				result = append(result, scanner.Bytes()...)
				result = append(result, br)
				continue
			}
		}
	}
	return result, nil
}

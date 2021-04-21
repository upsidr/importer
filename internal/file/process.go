package file

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// ProcessAnnotations reads annotations and generates ContentAfter.
func (f *File) ProcessAnnotations() error {
	result := []byte{}
	br := byte('\n')
	for line, data := range f.ContentPurged {
		result = append(result, data...)
		result = append(result, br)

		if a, found := f.Annotations[line+1]; found {
			processed, err := processSingleAnnotation(result, f.FileName, a)
			if err != nil {
				fmt.Printf("warning: %s", err)
				continue
			}
			result = processed
		}
	}
	f.ContentAfter = result
	return nil
}

func processSingleAnnotation(result []byte, filePath string, annotation *Annotation) ([]byte, error) {
	br := byte('\n')
	// Make sure the files are read based on the relative path
	dir := filepath.Dir(filePath)
	targetPath := dir + "/" + annotation.TargetPath
	file, err := os.Open(targetPath)
	if err != nil {
		// Purposely returning the byte slice as it contains data that were
		// populated prior to hitting this func
		return result, fmt.Errorf("could not open file '%s', skipping, %w", targetPath, err)
	}
	defer file.Close()

	// Prep
	reExport := regexp.MustCompile(ExportMarkerMarkdown)
	withinExportMarker := false
	currentLine := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++

		// Handle export marker imports
		match := reExport.FindStringSubmatch(scanner.Text())
		if len(match) != 0 {
			// match[1] is export_marker_name
			if match[1] == annotation.TargetExportMarker {
				withinExportMarker = true
			}
			// match[2] is exporter_marker_condition
			if match[2] == "end" {
				withinExportMarker = false
			}
			continue
		}
		if withinExportMarker {
			result = append(result, scanner.Bytes()...)
			result = append(result, br)
			continue
		}

		// Handle line number imports
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
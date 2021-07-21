package file

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const br = byte('\n')

// ProcessAnnotations reads annotations and generates ContentAfter.
//
// Internally, ContentAfter is generated from ContentPurged and Annotations.
// This walks thruogh each line of ContentPurged, and copies the data into a
// byte slice, and while processing each line, it checks if annotations are
// defined for the given line. If any annotation is registered, it would then
// process the target information to import.
//
// TODO: possibly remove error return, as it currently never returns any error.
func (f *File) ProcessAnnotations() error {
	result := []byte{}
	for line, data := range f.ContentPurged {
		result = append(result, data...)
		result = append(result, br)

		// Annotation is found for the given line. Before proceeding to the
		// next line, handle annotation and import the target data.
		if a, found := f.Annotations[line+1]; found {
			processed, err := processSingleAnnotation(f.FileName, a)
			if err != nil {
				fmt.Printf("warning: %s\n", err)
				continue
			}
			result = append(result, processed...)
		}
	}
	f.ContentAfter = result
	return nil
}

func processSingleAnnotation(filePath string, annotation *Annotation) ([]byte, error) {
	// TODO: Add support for URL https://github.com/upsidr/importer/issues/14

	// Make sure the files are read based on the relative path
	dir := filepath.Dir(filePath)
	targetPath := dir + "/" + annotation.TargetPath
	file, err := os.Open(targetPath)
	if err != nil {
		// Purposely returning the byte slice as it contains data that were
		// populated prior to hitting this func
		return nil, err
	}
	defer file.Close()

	fileType := filepath.Ext(filePath)
	switch fileType {
	case ".md":
		return processSingleAnnotationMarkdown(file, annotation)
	case ".yaml", ".yml":
		return processSingleAnnotationYAML(file, annotation)
	default:
		return processSingleAnnotationOther(file, annotation)
	}
}

func processSingleAnnotationMarkdown(file *os.File, annotation *Annotation) ([]byte, error) {
	result := []byte{}

	reExport := regexp.MustCompile(ExportMarkerMarkdown)
	withinExportMarker := false
	currentLine := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++

		// Find Exporter Marker
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

		// Handle export marker imports
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

func processSingleAnnotationYAML(file *os.File, annotation *Annotation) ([]byte, error) {
	result := []byte{}

	reExport := regexp.MustCompile(ExportMarkerYAML)
	withinExportMarker := false
	markerIndentation := 0
	currentLine := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++

		lineString := scanner.Text()
		lineData := scanner.Bytes()

		// Find Exporter Marker
		match := reExport.FindStringSubmatch(lineString)
		if len(match) != 0 {
			// match[1] is export_marker_indent
			if len(match[1]) > 0 {
				markerIndentation = len(match[1])
			}

			// match[2] is export_marker_name
			if match[2] == annotation.TargetExportMarker {
				withinExportMarker = true
			}
			// match[3] is exporter_marker_condition
			if match[3] == "end" {
				withinExportMarker = false
				markerIndentation = 0
			}
			continue
		}

		// Handle export marker imports
		if withinExportMarker {
			lineData = adjustIndentation(lineData, markerIndentation, annotation)
			result = append(result, lineData...)
			continue
		}

		// Handle line number imports
		if currentLine >= annotation.TargetLineFrom &&
			currentLine <= annotation.TargetLineTo {
			lineData = adjustIndentation(lineData, markerIndentation, annotation)
			result = append(result, lineData...)
			continue
		}
		for _, l := range annotation.TargetLines {
			if currentLine == l {
				lineData = adjustIndentation(lineData, markerIndentation, annotation)
				result = append(result, lineData...)
				continue
			}
		}
	}
	return result, nil
}

func processSingleAnnotationOther(file *os.File, annotation *Annotation) ([]byte, error) {
	result := []byte{}

	currentLine := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++

		// There is no export marker handling for unspecified file type. This
		// is because it is impossible to find what comment format is allowed
		// in the target file.
		// If there are specific use cases, it would need to be implemented
		// separately.

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

func adjustIndentation(lineData []byte, markerIndentation int, annotation *Annotation) []byte {
	lineString := string(lineData)

	// Check which indentation adjustment is used.
	// Absolute adjustment takes precedence over extra indentation.
	switch {
	case annotation.AbsoluteIndentation >= 0:
		actualIndent := len(lineString) - len(strings.TrimLeft(lineString, " "))
		switch {
		// Marker appears with more indentation than Absolute, and thus strip
		// extra indentations.
		case markerIndentation > annotation.AbsoluteIndentation:
			indentAdjustment := markerIndentation - annotation.AbsoluteIndentation
			lineData = lineData[indentAdjustment:]

		// Marker has less indentation than Absolute wants, and thus prepend
		// the indent diff.
		case markerIndentation < annotation.AbsoluteIndentation:
			indentAdjustment := annotation.AbsoluteIndentation - markerIndentation
			lineData = prependWhitespaces(lineData, indentAdjustment)
		case actualIndent < markerIndentation:
			// TODO: Handle case where indentation is less than marker indentation
		}
	case annotation.ExtraIndentation > 0:
		lineData = prependWhitespaces(lineData, annotation.ExtraIndentation)
	}
	lineData = append(lineData, br)
	return lineData
}

func prependWhitespaces(x []byte, count int) []byte {
	empty := bytes.Repeat([]byte(" "), count)
	// x = append(x, empty...)
	// copy(x[count:], x)
	// copy(x, empty)
	// return x

	return append(empty, x...)
}

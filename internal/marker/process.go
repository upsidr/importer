package marker

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ProcessMarker processes the marker data to generate the byte array of import
// target. Marker validation is assumed by using NewMarker.
//
// `importingFilePath` input is used for resolving relative filepath to find
// the import target.
func (m *Marker) ProcessMarker(importingFilePath string) ([]byte, error) {
	// TODO: Add support for URL https://github.com/upsidr/importer/issues/14

	// Make sure the files are read based on the relative path
	dir := filepath.Dir(importingFilePath)
	targetPath := dir + "/" + m.TargetPath
	file, err := os.Open(targetPath)
	if err != nil {
		// Purposely returning the byte slice as it contains data that were
		// populated prior to hitting this func
		return nil, err
	}
	defer file.Close()

	fileType := filepath.Ext(importingFilePath)
	switch fileType {
	case ".md":
		return m.processSingleMarkerMarkdown(file)
	case ".yaml", ".yml":
		return m.processSingleMarkerYAML(file)
	default:
		return m.processSingleMarkerOther(file)
	}
}

const br = byte('\n')

func (m *Marker) processSingleMarkerMarkdown(file *os.File) ([]byte, error) {
	result := []byte{}

	reExport := regexp.MustCompile(ExporterMarkerMarkdown)
	withinExportMarker := false
	currentLine := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++

		// Find Exporter Marker
		match := reExport.FindStringSubmatch(scanner.Text())
		if len(match) != 0 {
			// match[1] is export_marker_name
			if match[1] == m.TargetExportMarker {
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
		if currentLine >= m.TargetLineFrom &&
			currentLine <= m.TargetLineTo {
			result = append(result, scanner.Bytes()...)
			result = append(result, br)
			continue
		}
		for _, l := range m.TargetLines {
			if currentLine == l {
				result = append(result, scanner.Bytes()...)
				result = append(result, br)
				continue
			}
		}
	}
	return result, nil
}

func (m *Marker) processSingleMarkerYAML(file *os.File) ([]byte, error) {
	result := []byte{}

	reExport := regexp.MustCompile(ExporterMarkerYAML)
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
			if match[2] == m.TargetExportMarker {
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
			lineData = adjustIndentation(lineData, markerIndentation, m)
			result = append(result, lineData...)
			continue
		}

		// Handle line number imports
		if currentLine >= m.TargetLineFrom &&
			currentLine <= m.TargetLineTo {
			lineData = adjustIndentation(lineData, markerIndentation, m)
			result = append(result, lineData...)
			continue
		}
		for _, l := range m.TargetLines {
			if currentLine == l {
				lineData = adjustIndentation(lineData, markerIndentation, m)
				result = append(result, lineData...)
				continue
			}
		}
	}
	return result, nil
}

func (m *Marker) processSingleMarkerOther(file *os.File) ([]byte, error) {
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
		if currentLine >= m.TargetLineFrom &&
			currentLine <= m.TargetLineTo {
			result = append(result, scanner.Bytes()...)
			result = append(result, br)
			continue
		}
		for _, l := range m.TargetLines {
			if currentLine == l {
				result = append(result, scanner.Bytes()...)
				result = append(result, br)
				continue
			}
		}
	}
	return result, nil
}

func adjustIndentation(lineData []byte, markerIndentation int, marker *Marker) []byte {
	// If no indentation setup is done, simply return as is
	if marker.Indentation == nil {
		lineData = append(lineData, br)
		return lineData
	}

	lineString := string(lineData)
	indentLength := marker.Indentation.Length
	// Check which indentation adjustment is used.
	// Absolute adjustment takes precedence over extra indentation.
	switch marker.Indentation.Mode {
	case AbsoluteIndentation:
		actualIndent := len(lineString) - len(strings.TrimLeft(lineString, " "))
		switch {
		// Marker appears with more indentation than Absolute, and thus strip
		// extra indentations.
		case markerIndentation > indentLength:
			indentAdjustment := markerIndentation - indentLength
			lineData = lineData[indentAdjustment:]

		// Marker has less indentation than Absolute wants, and thus prepend
		// the indent diff.
		case markerIndentation < indentLength:
			indentAdjustment := indentLength - markerIndentation
			lineData = prependWhitespaces(lineData, indentAdjustment)
		case actualIndent < markerIndentation:
			// TODO: Handle case where indentation is less than marker indentation
		}
	case ExtraIndentation:
		lineData = prependWhitespaces(lineData, indentLength)
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

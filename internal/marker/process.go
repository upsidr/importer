package marker

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/upsidr/importer/internal/regexpplus"
)

// ProcessMarkerData processes the marker data to generate the byte array of
// import target. Marker validation is assumed by using NewMarker.
//
// `importingFilePath` input is used for resolving relative filepath to find
// the import target.
func (m *Marker) ProcessMarkerData(importingFilePath string) ([]byte, error) {
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

	// reExport := regexp.MustCompile(ExporterMarkerYAML)
	isNested := false
	nestedUnder := ""
	exporterMarkerIndentation := 0
	currentLine := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++

		lineString := scanner.Text()
		lineData := scanner.Bytes()

		switch {
		// Handle line number range
		case m.TargetLineFrom > 0:
			if currentLine >= m.TargetLineFrom &&
				currentLine <= m.TargetLineTo {
				lineData = append(lineData, br)
				result = append(result, lineData...)
				continue
			}

		// Handle line number slice
		case len(m.TargetLines) > 0:
			for _, l := range m.TargetLines {
				if currentLine == l {
					lineData = append(lineData, br)
					result = append(result, lineData...)
					continue
				}
			}

		// Handle ExporterMarker
		case m.TargetExportMarker != "":
			// Find Exporter Marker
			matches, err := regexpplus.MapWithNamedSubgroups(lineString, ExporterMarkerYAML)
			if errors.Is(err, regexpplus.ErrNoMatch) {
				// This line is not Exporter Marker. If there has been some
				// marker found already, append the line and continue
				if isNested {
					lineData = adjustIndentation(lineData, exporterMarkerIndentation, m.Indentation)
					result = append(result, lineData...)
				}
				continue
			}
			if err != nil {
				panic(err) // Unknown error, should not happen
			}

			if exporterName, found := matches["export_marker_name"]; found {
				// Ignore unrelated marker
				if exporterName != m.TargetExportMarker || exporterName == "" {
					continue
				}

				if x, found := matches["export_marker_indent"]; found {
					indent := len(x) - len(strings.TrimLeft(x, " "))
					exporterMarkerIndentation = indent
				}

				if exporterCondition, found := matches["exporter_marker_condition"]; found {
					if exporterName == nestedUnder && exporterCondition == "end" {
						isNested = false
						exporterMarkerIndentation = 0
						continue
					}
				}
				nestedUnder = exporterName
				isNested = true
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

func adjustIndentation(lineData []byte, exporterMarkerIndent int, importerIndentation *Indentation) []byte {
	// If no indentation setup is done, simply return as is
	if importerIndentation == nil {
		lineData = append(lineData, br)
		return lineData
	}

	// Check which indentation adjustment is used.
	// Absolute adjustment takes precedence over extra indentation.
	switch importerIndentation.Mode {
	case AbsoluteIndentation:
		lineData = handleAbsoluteIndentation(lineData, exporterMarkerIndent, importerIndentation.Length)
	case ExtraIndentation:
		lineData = prependWhitespaces(lineData, importerIndentation.Length)
	case AlignIndentation:
		lineData = handleAbsoluteIndentation(lineData, exporterMarkerIndent, importerIndentation.MarkerIndentation)
	}
	lineData = append(lineData, br)
	return lineData
}

// handleAbsoluteIndentation updates the lineData with provided indentation.
//
// There are 3 different indent information to handle:
//
//   A. Indent target in absolute value
//   B. Original preceding indent in the lineData
//   C. Indent of Exporter Marker
//
// With the absolute indentation handling, A and B are obviously required. For
// C, it is important that YAML tree structure is persisted, and thus it is
// crucial to know how much indent Exporter Marker has.
//
// With the above covered, there are a few cases to handle:
//
//   1. Original lineData has more indentation than target indentation
//   2. Original lineData has fewer indentation than target indentation
//   3. Original lineData has fewer preceding indentation than Importer Marker
//
// For the Case 1. and 2., the diff needs to be calculated to ensure correct
// indentation.
//
// For the Case 3., it's not clear what we should expect. This is currently not
// handled, and it may need to be an error.
func handleAbsoluteIndentation(lineData []byte, exportMarkerIndent, targetIndent int) []byte {
	lineString := string(lineData)
	currenttIndent := len(lineString) - len(strings.TrimLeft(lineString, " "))

	switch {
	// Case 1.
	// Marker appears with more indentation than Absolute, and thus strip
	// extra indentations.
	case exportMarkerIndent >= targetIndent:
		indentAdjustment := exportMarkerIndent - targetIndent
		return lineData[indentAdjustment:]

	// Case 2.
	// Marker has less indentation than Absolute wants, and thus prepend
	// the indent diff.
	case exportMarkerIndent < targetIndent:
		indentAdjustment := targetIndent - exportMarkerIndent
		return prependWhitespaces(lineData, indentAdjustment)

	// Case 3.
	case currenttIndent < exportMarkerIndent:
		// TODO: Handle case where indentation is less than marker indentation
	}
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

package marker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/upsidr/importer/internal/regexpplus"
)

// ProcessMarkerData processes the marker data to generate the byte array of
// import target. Marker validation is assumed by using NewMarker.
//
// `importingFilePath` input is used for resolving relative filepath to find
// the import target.
func (m *Marker) ProcessMarkerData(importingFilePath string) ([]byte, error) {
	var file io.Reader

	targetFile := m.ImportTargetFile.File
	switch m.ImportTargetFile.Type {
	case PathBased:
		// Make sure the files are read based on the relative path
		dir := filepath.Dir(targetFile)
		targetPath := dir + "/" + m.ImportTargetFile.File
		f, err := os.Open(targetPath)
		if err != nil {
			// TODO: This note is no longer true - need to review what it was meant to be.
			// Purposely returning the byte slice as it contains data that were
			// populated prior to hitting this func
			return nil, err
		}
		defer f.Close()
		file = f
	case URLBased:
		u, err := preprocessURL(targetFile)
		if err != nil {
			return nil, fmt.Errorf("%w of '%s'", ErrInvalidURL, targetFile)
		}
		r, err := http.Get(u)
		if err != nil {
			return nil, fmt.Errorf("%w, %v", ErrGetMarkerTarget, err)
		}
		if r.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%w of '%s'", ErrNonSuccessCode, r.Status)
		}
		defer r.Body.Close()
		file = r.Body
	default:
		return nil, fmt.Errorf("%w", ErrNoFileInput)
	}

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

func (m *Marker) processSingleMarkerMarkdown(file io.Reader) ([]byte, error) {
	result := []byte{}

	withinExportMarker := false
	currentLine := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		currentLine++

		// Find Exporter Marker
		matches, err := regexpplus.MapWithNamedSubgroups(scanner.Text(), ExporterMarkerMarkdown)
		if err != nil && !errors.Is(err, regexpplus.ErrNoMatch) {
			panic(err) // Unknown error, should not happen
		}

		if len(matches) != 0 {
			if exporterName, found := matches["export_marker_name"]; found &&
				exporterName == m.ImportLogic.ExporterMarker {
				withinExportMarker = true
			}
			if exporterCondition, found := matches["exporter_marker_condition"]; found &&
				exporterCondition == "end" {
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
		if currentLine >= m.ImportLogic.LineFrom &&
			currentLine <= m.ImportLogic.LineTo {
			result = append(result, scanner.Bytes()...)
			result = append(result, br)
			continue
		}
		for _, l := range m.ImportLogic.Lines {
			if currentLine == l {
				result = append(result, scanner.Bytes()...)
				result = append(result, br)
				continue
			}
		}
	}
	return result, nil
}

func (m *Marker) processSingleMarkerYAML(file io.Reader) ([]byte, error) {
	result := []byte{}

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
		case m.ImportLogic.LineFrom > 0:
			if currentLine >= m.ImportLogic.LineFrom &&
				currentLine <= m.ImportLogic.LineTo {
				lineData = append(lineData, br)
				result = append(result, lineData...)
				continue
			}

		// Handle line number slice
		case len(m.ImportLogic.Lines) > 0:
			for _, l := range m.ImportLogic.Lines {
				if currentLine == l {
					lineData = append(lineData, br)
					result = append(result, lineData...)
					continue
				}
			}

		// Handle ExporterMarker
		case m.ImportLogic.ExporterMarker != "":
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
				if exporterName != m.ImportLogic.ExporterMarker || exporterName == "" {
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

func (m *Marker) processSingleMarkerOther(file io.Reader) ([]byte, error) {
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
		if currentLine >= m.ImportLogic.LineFrom &&
			currentLine <= m.ImportLogic.LineTo {
			result = append(result, scanner.Bytes()...)
			result = append(result, br)
			continue
		}
		for _, l := range m.ImportLogic.Lines {
			if currentLine == l {
				result = append(result, scanner.Bytes()...)
				result = append(result, br)
				continue
			}
		}
	}
	return result, nil
}

// preprocessURL checks the incoming address input and returns the updated
// URL.
//
// This includes updates such as github.com address to be based on
// raw.githubusercontent.com for easier reference.
func preprocessURL(address string) (string, error) {
	u, err := url.ParseRequestURI(address)
	if err != nil {
		return "", fmt.Errorf("%w, %v", ErrInvalidPath, err)
	}

	// For non-github.com address, simply return as is
	if u.Host != "github.com" {
		return address, nil
	}

	// TODO: naïve implementation, this won't work for username "blob" or repo name "blob"
	ps := strings.SplitN(u.Path, "/blob/", 2)
	if len(ps) == 1 {
		return address, nil
	}

	u.Host = "raw.githubusercontent.com"
	u.Path = strings.Join(ps, "/")

	return u.String(), nil
}

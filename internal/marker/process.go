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
	var file io.Reader
	switch {
	case m.TargetPath != "":
		// Make sure the files are read based on the relative path
		dir := filepath.Dir(importingFilePath)
		targetPath := dir + "/" + m.TargetPath
		f, err := os.Open(targetPath)
		if err != nil {
			// Purposely returning the byte slice as it contains data that were
			// populated prior to hitting this func
			return nil, err
		}
		defer f.Close()
		file = f
	case m.TargetURL != "":
		u, err := preprocessURL(m.TargetURL)
		if err != nil {
			return nil, err
		}
		r, err := http.Get(u)
		if err != nil {
			return nil, err // TODO: test coverage
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

func (m *Marker) processSingleMarkerYAML(file io.Reader) ([]byte, error) {
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

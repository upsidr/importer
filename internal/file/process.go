package file

import (
	"bufio"
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/upsidr/importer/internal/marker"
)

const br = byte('\n')

// ProcessMarkers reads markers and generates ContentAfter.
//
// Internally, ContentAfter is generated from ContentPurged and Markers.
// This walks thruogh each line of ContentPurged, and copies the data into a
// byte slice, and while processing each line, it checks if markers are
// defined for the given line. If any marker is registered, it would then
// process the target information to import.
//
// TODO: possibly remove error return, as it currently never returns any error.
func (f *File) ProcessMarkers() error {
	result := []byte{}
	for line, data := range f.ContentPurged {
		result = append(result, data...)
		result = append(result, br)

		// Marker is found for the given line. Before proceeding to the
		// next line, handle marker and import the target data.
		if marker, found := f.Markers[line+1]; found {
			processed, err := marker.ProcessMarkerData(f.FileName)
			if err != nil {
				fmt.Printf("Warning: error while processing '%s': %s\n", marker.Name, err)
				continue
			}
			result = append(result, processed...)
		}
	}
	f.ContentAfter = result
	return nil
}

// RemoveMarkers removes Importer markers. This is useful for generated files
// to have no marker input.
func (f *File) RemoveMarkers() {
	fileType := filepath.Ext(f.FileName)

	var importerRe *regexp.Regexp
	switch fileType {
	case ".md":
		importerRe = regexp.MustCompile(marker.ImporterMarkerMarkdown)
	case ".yaml", ".yml":
		importerRe = regexp.MustCompile(marker.ImporterMarkerYAML)
	default:
		// File that does not have supporting marker setup will be simply
		// ignored.
	}

	var exporterRe *regexp.Regexp
	switch fileType {
	case ".md":
		exporterRe = regexp.MustCompile(marker.ExporterMarkerMarkdown)
	case ".yaml", ".yml":
		exporterRe = regexp.MustCompile(marker.ExporterMarkerYAML)
	default:
		// File that does not have supporting marker setup will be simply
		// ignored.
	}

	if importerRe == nil || exporterRe == nil {
		return
	}

	newResult := []byte{}

	scanner := bufio.NewScanner(bytes.NewReader(f.ContentAfter))
	for scanner.Scan() {
		currentLine := scanner.Bytes()
		if s := importerRe.Find(currentLine); s != nil {
			// Importer Marker found, ignore
			continue
		}
		if s := exporterRe.Find(currentLine); s != nil {
			// Exporter Marker found, ignore
			continue
		}

		currentLine = append(currentLine, []byte("\n")...)
		newResult = append(newResult, currentLine...)
	}

	f.ContentAfter = newResult

	return
}

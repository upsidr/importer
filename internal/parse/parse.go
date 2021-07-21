package parse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"

	"github.com/upsidr/importer/internal/file"
)

var (
	ErrUnsupportedFileType   = errors.New("unsupported file type")
	ErrNoInput               = errors.New("no file content found")
	ErrInvalidPath           = errors.New("invalid path provided")
	ErrInvalidSyntax         = errors.New("invalid syntax given")
	ErrNoMatchingAnnotations = errors.New("no matching annotations found, annotation must be a begin/end pair")
)

// File holds onto file data.
type File struct {
	FileName string

	// ContentBefore holds the file content as it was before processing. The
	// first slice represents the line number, and the second is for the actual
	// data.
	ContentBefore [][]byte

	// ContentPurged holds the file content, but removes the parts between
	// importer annotation begin/end. The first slice represents the line
	// number, and the second is for the actual data.
	ContentPurged [][]byte

	// ContentAfter holds the file content after the import has been run. This
	// only holds the actual data in byte slice representation.
	ContentAfter []byte

	// Annotations is an array holding onto each annotation block.
	Annotations map[int]file.Annotation
}

// Parse reads filename and input, and parses data in the file.
//
// The steps are as follows:
//
// 	1. Read input data
// 	2. Scan each line
// 	3. Look for regex match for annotation
// 	4. Save matched line number and options found
// 	5. Verify parsed data, and return
//
// If any of the above steps failed, it would return an error. This function
// does not populate the ContentAfter.
func Parse(fileName string, input io.Reader) (*file.File, error) {
	if input == nil {
		return nil, ErrNoInput
	}

	fileType := filepath.Ext(fileName)

	switch fileType {
	case ".md":
		return parse(ImporterAnnotationMarkdown, fileName, input)
	case ".yaml", ".yml":
		return parse(ImporterAnnotationYAML, fileName, input)
	default:
		return nil, fmt.Errorf("%w, '%s' provided", ErrUnsupportedFileType, fileType)
	}
}

// parse reads file input using scanner. This reads the input line by line, and
// store the data into File data. Parsing the data stores 3 sets of data: file
// content as is, annotation details, and file content with all data between
// annotation pairs purged.
func parse(annotation string, fileName string, input io.Reader) (*file.File, error) {
	f := &file.File{FileName: fileName}
	re := regexp.MustCompile(annotation)

	annotations := map[int]*file.Annotation{}
	matches := map[string]matchHolder{}

	currentLine := 0
	inNested := false // Flag to check if the data is between annotations
	nestedUnder := "" // Name to check for annotation pair ending

	// NOTE:
	// For *File.contentXyz, I'm purposely making the first item in slice empty
	// for readability. This shouldn't be necessary, but with this approach,
	// the slice index matches the line number, and is easy to get my head
	// around for now.
	f.ContentBefore = make([]string, 0)
	f.ContentPurged = make([]string, 0)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		currentLine++
		currentStr := scanner.Text()
		f.ContentBefore = append(f.ContentBefore, currentStr)

		// Look for marker match
		match := re.FindStringSubmatch(currentStr)
		if len(match) == 0 {
			// If the line appears within some other marker set, remove the line.
			if inNested {
				continue
			}
			// Otherwise ensure the marker itself is a part of purged data.
			f.ContentPurged = append(f.ContentPurged, currentStr)

			// There is no further action needed for matched line, and thus continue.
			continue
		}

		// Parse regex match into groups to handle annotation
		subgroupName := match[1] // regex 1st subgroup. Index 0 is for full string.

		// Ensure this is the top most annotation. If a nested annotation is
		// found within another annotation, ignore it. This is because we
		// should be handling those nested annotations in those target files
		// instead.
		// This means, in any file in question, the parse logic only looks at
		// one file and its direct dependencies.
		// TODO: Handle file dependencies with AST
		if inNested && nestedUnder != subgroupName {
			continue
		}

		nestedUnder = subgroupName

		// At this point, the annotation is important, and we need to process
		// the line further.
		// Note that, ContentPurged does not contain any data that's wrapped
		// between annotations. Those lines will be kept as an empty byte slice
		// for further processing later to create ContentAfter.
		f.ContentPurged = append(f.ContentPurged, currentStr)

		// Annotations must match up to create a pair. If it isn't a proper
		// pair, it is treated as broken. For that reason, we need to keep
		// track of already found match.
		matchData := matchHolder{}
		if data, found := matches[subgroupName]; found {
			// TODO: Handle case where the same subgroup name gets used multiple times.
			matchData = data
		}

		for i, n := range re.SubexpNames() {
			matchedContent := match[i]
			switch n {
			// The first subgroup is the name, which is used as the map key.
			case "importer_name":
				continue
			case "importer_marker":
				if matchedContent == "begin" {
					inNested = true
					matchData.isBeginFound = true
					matchData.lineToInsertAt = len(f.ContentPurged)
				}
				if matchedContent == "end" {
					inNested = false
					nestedUnder = ""
					matchData.isEndFound = true
				}
			case "importer_option":
				if matchedContent != "" { // TODO: skipping empty string like this as end annotation shouldn't override
					matchData.options = matchedContent
				}
			}
		}
		matches[subgroupName] = matchData
	}

	for name, data := range matches {
		annotation, err := convert(name, data)
		if err != nil {
			// TODO: err should be handled rather than simply ignored.
			//       This is fine for now as error is used for internal logic
			//       only, but shuold be fixed.
			continue
		}

		annotations[annotation.LineToInsertAt] = annotation
	}

	f.Annotations = annotations

	return f, nil
}

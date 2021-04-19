package parse

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"

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
// 	4. Save matched line number and option found
// 	5. Verify parsed data, and return
//
// If any of the above steps failed, it would return an error. This function
// does not populate the ContentAfter.
func Parse(fileName string, input io.Reader) (*file.File, error) {
	if input == nil {
		return nil, ErrNoInput
	}

	fileType := filepath.Ext(fileName)

	var result *file.File
	var err error
	switch fileType {
	case ".md":
		result, err = parseMarkdown(fileName, input)
		if err != nil {
			return nil, err // TODO: test coverage
		}
	case ".yaml", ".yml":
		fmt.Printf("yaml") // TODO: implement
	default:
		return nil, fmt.Errorf("%w, '%s' provided", ErrUnsupportedFileType, fileType)
	}

	return result, nil
}

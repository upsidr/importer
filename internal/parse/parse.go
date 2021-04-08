package parse

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/upsidr/importer/internal/file"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrNoInput             = errors.New("no file content found")
)

// File holds onto file data.
type File struct {
	FileName string

	// fileType is derived from FileName, which is simply represented using
	// extension format.
	fileType string

	// ContentBefore holds the file content as it was before processing. The
	// first slice represents the line number, and the second is for the actual
	// data.
	ContentBefore [][]byte

	// ContentPurged holds the file coontent, but removes the parts between
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
// If any of the above steps failed, it would return an error.
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

	err = result.UpdateWithAnnotations()
	if err != nil {
		return nil, err // TODO: test coverage
	}

	return result, nil
}

// matchHolder is a temporary data holder, which is used to ensure validity of
// annotation data.
type matchHolder struct {
	isBeginFound   bool
	isEndFound     bool
	lineToInsertAt int
	options        []string
}

func convert(name string, match matchHolder) (*file.Annotation, error) {
	if !match.isBeginFound || !match.isEndFound {
		return nil, errors.New("no matching annotations found, annotation must be a begin/end pair")
	}
	reImportTarget := regexp.MustCompile(OptionFilePathIndicator)

	result := &file.Annotation{
		Name:           name,
		LineToInsertAt: match.lineToInsertAt,
	}

	for _, opt := range match.options {
		match := reImportTarget.FindAllStringSubmatch(opt, -1)
		if len(match) == 0 {
			continue
		}

		for _, ms := range match {
			for i, n := range reImportTarget.SubexpNames() {
				matchedContent := ms[i]
				switch n {
				case "importer_target_path":
					result.TargetPath = matchedContent
				case "importer_target_lines":
					lines := []int{}

					// Handle case where the input is something like "6~22"
					if strings.Contains(matchedContent, "~") {
						ls := strings.Split(matchedContent, "~")
						lb, err := strconv.Atoi(ls[0])
						if err != nil {
							return nil, fmt.Errorf("error: found non-numeric input for line number lower bound, %v", err)
						}
						ub, err := strconv.Atoi(ls[1])
						if err != nil {
							return nil, fmt.Errorf("error: found non-numeric input for line number upper bound, %v", err)
						}
						// Add line numbers to the slice.
						// This way, we can support comma separated list, etc.
						for i := lb; i <= ub; i++ {
							lines = append(lines, i)
						}
					}
					result.TargetLines = lines
				}
			}
		}
	}

	return result, nil
}

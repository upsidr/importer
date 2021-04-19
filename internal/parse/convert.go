package parse

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/upsidr/importer/internal/file"
)

// matchHolder is a temporary data holder, which is used to ensure validity of
// annotation data.
type matchHolder struct {
	isBeginFound   bool
	isEndFound     bool
	lineToInsertAt int
	options        string
}

func convert(name string, match matchHolder) (*file.Annotation, error) {
	if !match.isBeginFound || !match.isEndFound {
		return nil, fmt.Errorf("%w", ErrNoMatchingAnnotations)
	}
	reImportTarget := regexp.MustCompile(OptionFilePathIndicator)

	result := &file.Annotation{
		Name:           name,
		LineToInsertAt: match.lineToInsertAt,
	}

	ms := reImportTarget.FindAllStringSubmatch(match.options, -1)

	// No option provided
	// TODO: Handle this case better
	if len(ms) == 0 {
		return result, nil
	}

	for _, m := range ms {
		for i, n := range reImportTarget.SubexpNames() {
			matchedContent := m[i]
			switch n {
			case "importer_target_path":
				if err := processTargetPath(result, matchedContent); err != nil {
					return nil, err // TODO: Add test coverage when more validation is added to path check
				}
			case "importer_target_detail":
				if err := processTargetDetail(result, matchedContent); err != nil {
					return nil, err
				}
			}
		}
	}

	return result, nil
}

// processTargetPath processes string input of import target path.
//
// Target path can be 2 forms.
//   - Relative or absolute path to local file
//   - URL to retrieve the file from
//
// TODO: URL handling to be supported
func processTargetPath(annotation *file.Annotation, input string) error {
	// TODO: Add more validation
	if input == "" {
		return fmt.Errorf("%w", ErrInvalidPath)
	}

	annotation.TargetPath = input

	return nil
}

// processTargetDetail processes string input of import detail, which contains
// some detail of what to import from the target.
//
// Target detail can be in various forms.
//   - Export marker, e.g. "[some_export_marker]", where it looks for
//     "some_export_marker" within the target file. This can hold comma
//     separated entries.
//   - Line range, e.g. "6~22" meaning line 6 to 22.
//   - Open line range, e.g. "~22" for line 1 to 22, "6~" for line 6 to end of
//     file.
//   - Line selection, e.g. "1,5,7" meaning line 1, 5 and 7.
func processTargetDetail(annotation *file.Annotation, input string) error {
	exportMarker := regexp.MustCompile(`\[(\S+)\]`)

	marker := exportMarker.FindStringSubmatch(input)
	switch {
	// Handle export marker
	case marker != nil:
		annotation.TargetExportMarker = string(marker[1])

	// Handle line range annotation with commas
	case strings.Contains(input, ","):
		targetLines := []int{}

		// Handle comma separated numbers
		nums := strings.Split(input, ",")
		for _, num := range nums {
			// Handle tilde based range notation
			if strings.Contains(num, "~") {
				ls := strings.Split(num, "~")

				// if conversion fails, simply ignore to try processing the rest
				lowerBound, _ := strconv.Atoi(ls[0])
				upperBound, _ := strconv.Atoi(ls[1])

				// Add line numbers to the slice.
				// This way, we can support comma separated list, etc.
				for i := lowerBound; i <= upperBound; i++ {
					targetLines = append(targetLines, i)
				}
			}

			// Handle single number
			lineNumber, err := strconv.Atoi(num)
			if err != nil {
				continue
			}
			targetLines = append(targetLines, lineNumber)
		}

		annotation.TargetLines = targetLines

	// Handle single line range
	case strings.Contains(input, "~"):
		ls := strings.Split(input, "~")
		lb := ls[0]
		ub := ls[1]

		if lb != "" {
			lowerBound, err := strconv.Atoi(lb)
			if err != nil {
				return fmt.Errorf("%w, %v", ErrInvalidSyntax, err)
			}
			annotation.TargetLineFrom = lowerBound
		}

		if ub != "" {
			upperBound, err := strconv.Atoi(ub)
			if err != nil {
				return fmt.Errorf("%w, %v", ErrInvalidSyntax, err)
			}
			annotation.TargetLineTo = upperBound
		}

	default:
		i, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("%w, %v", ErrInvalidSyntax, err)
		}
		annotation.TargetLines = append(annotation.TargetLines, i)
	}

	return nil
}

package parse

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/upsidr/importer/internal/file"
	"github.com/upsidr/importer/internal/regexpplus"
)

// matchHolder is a temporary data holder, which is used to ensure validity of
// marker data.
type matchHolder struct {
	isBeginFound   bool
	isEndFound     bool
	lineToInsertAt int
	options        string
}

func processMarker(name string, match matchHolder) (*file.Marker, error) {
	if !match.isBeginFound || !match.isEndFound {
		return nil, fmt.Errorf("%w", ErrNoMatchingMarker)
	}

	result := &file.Marker{
		Name:           name,
		LineToInsertAt: match.lineToInsertAt,
	}

	err := processFileOption(result, match)
	if err != nil {
		return nil, err
	}

	err = processIndentOption(result, match)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func processFileOption(marker *file.Marker, match matchHolder) error {
	matches, err := regexpplus.MapWithNamedSubgroups(match.options, OptionFilePathIndicator)
	if err != nil {
		return nil // Indent options are not required, and thus simply ignore if no match
	}

	if targetPath, ok := matches["importer_target_path"]; ok {
		if err := processTargetPath(marker, targetPath); err != nil {
			return err // TODO: Add test coverage when more validation is added to path check
		}
	}
	if targetDetail, ok := matches["importer_target_detail"]; ok {
		if err := processTargetDetail(marker, targetDetail); err != nil {
			return err
		}
	}

	return nil
}

func processIndentOption(marker *file.Marker, match matchHolder) error {
	matches, err := regexpplus.MapWithNamedSubgroups(match.options, OptionIndentMode)
	if err != nil {
		return nil // Indent options are not required, and thus simply ignore if no match
	}

	if indent, ok := matches["importer_indent_mode"]; ok {
		switch indent {
		case "absolute":
			marker.Indentation = &file.Indentation{Mode: file.AbsoluteIndentation}
		case "extra":
			marker.Indentation = &file.Indentation{Mode: file.ExtraIndentation}
		default:
			return errors.New("unsupported indentation mode")
		}
	}

	if lengthInput, ok := matches["importer_indent_length"]; ok {
		// Indentation length can be handled only when indentation mode
		// is specified. As RegEx handling should start from mode handling,
		// marker.Indentation shouldn't be nil at this point.

		length, err := strconv.Atoi(lengthInput)
		if err != nil {
			return err
		}
		marker.Indentation.Length = length
	}

	return nil
}

// processTargetPath processes string input of import target path.
//
// Target path can be 2 forms.
//   - Relative or absolute path to local file
//   - URL to retrieve the file from
//
// TODO: URL handling to be supported
func processTargetPath(marker *file.Marker, input string) error {
	// TODO: Add more validation
	if input == "" {
		return fmt.Errorf("%w", ErrInvalidPath)
	}

	marker.TargetPath = input

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
func processTargetDetail(marker *file.Marker, input string) error {
	exportMarker := regexp.MustCompile(`\[(\S+)\]`)

	markerRegex := exportMarker.FindStringSubmatch(input)
	switch {
	// Handle export marker
	case markerRegex != nil:
		marker.TargetExportMarker = string(markerRegex[1])

	// Handle line range marker with commas
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

		marker.TargetLines = targetLines

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
			marker.TargetLineFrom = lowerBound
		}

		marker.TargetLineTo = math.MaxInt32 // TODO: Consider making this Int64
		if ub != "" {
			upperBound, err := strconv.Atoi(ub)
			if err != nil {
				return fmt.Errorf("%w, %v", ErrInvalidSyntax, err)
			}
			marker.TargetLineTo = upperBound
		}

	default:
		i, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("%w, %v", ErrInvalidSyntax, err)
		}
		marker.TargetLines = append(marker.TargetLines, i)
	}

	return nil
}

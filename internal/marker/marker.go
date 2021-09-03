package marker

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/upsidr/importer/internal/regexpplus"
)

// Marker holds on to the data required for importer processing. This does
// not hold the target file content itself, and that needs to be handled
// separately.
//
// NewMarker function ensures the data validity of Marker for further
// processing.
type Marker struct {
	Name           string
	LineToInsertAt int

	TargetPath string
	TargetURL  string

	TargetExportMarker string
	TargetLines        []int
	TargetLineFrom     int
	TargetLineTo       int

	Indentation *Indentation

	// TODO: Add insert style such as code verbatim, details, quotes, etc.
}

type IndentationMode int

const (
	// Reserve 0 value as invalid
	AbsoluteIndentation IndentationMode = iota + 1
	ExtraIndentation
	AlignIndentation
	KeepIndentation
)

// Indentation holds additional indentation handling option.
type Indentation struct {
	Mode              IndentationMode
	Length            int
	MarkerIndentation int
}

func NewMarker(raw *RawMarker) (*Marker, error) {
	err := raw.Validate()
	if err != nil {
		return nil, err
	}

	result := &Marker{
		Name:           raw.Name,
		LineToInsertAt: raw.LineToInsertAt,
	}

	err = processFileOption(result, raw)
	if err != nil {
		return nil, err
	}

	err = processIndentOption(result, raw)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func processFileOption(marker *Marker, match *RawMarker) error {
	matches, err := regexpplus.MapWithNamedSubgroups(match.Options, OptionFilePathIndicator)
	if err != nil {
		return fmt.Errorf("%w, import target option is missing", ErrInvalidSyntax)
	}

	if targetPath, found := matches["importer_target_path"]; found {
		if err := processTargetPath(marker, targetPath); err != nil {
			return err
		}
	}
	if targetDetail, found := matches["importer_target_detail"]; found {
		if err := processTargetDetail(marker, targetDetail); err != nil {
			return err
		}
	}

	return nil
}

func processIndentOption(marker *Marker, match *RawMarker) error {
	matches, err := regexpplus.MapWithNamedSubgroups(match.Options, OptionIndentMode)
	if err != nil {
		return nil // Indent options are not required, and thus simply ignore if no match
	}

	if indent, found := matches["importer_indent_mode"]; found {
		switch indent {
		case "absolute":
			marker.Indentation = &Indentation{Mode: AbsoluteIndentation}
		case "extra":
			marker.Indentation = &Indentation{Mode: ExtraIndentation}
		case "align":
			markerIndentation := len(match.PrecedingIndentation) - len(strings.TrimLeft(match.PrecedingIndentation, " ")) // Naïve count
			marker.Indentation = &Indentation{
				Mode:              AlignIndentation,
				MarkerIndentation: markerIndentation,
			}
			return nil // Align option does not care length information
		default:
			return errors.New("unsupported indentation mode") // This shouldn't happen with the underlying regex
		}
	}

	if lengthInput, found := matches["importer_indent_length"]; found {
		// Indentation length can be handled only when indentation mode
		// is specified. As RegEx handling should start from mode handling,
		// marker.Indentation shouldn't be nil at this point.

		length, err := strconv.Atoi(lengthInput)
		if err != nil {
			return fmt.Errorf("%w, %v", ErrInvalidSyntax, err)
		}
		marker.Indentation.Length = length
	}

	return nil
}

// processTargetPath processes string input of import target path.
//
// Target path can be 2 forms.
//   - URL to retrieve the file from
//   - Relative or absolute path to local file
func processTargetPath(marker *Marker, input string) error {
	switch {
	// TODO: Naïve implementation, fix this
	case strings.HasPrefix(input, "http://"),
		strings.HasPrefix(input, "https://"):
		_, err := url.ParseRequestURI(input)
		if err != nil {
			return err
		}
		marker.TargetURL = input
	default:
		_, file := filepath.Split(input)
		if file == "" {
			return fmt.Errorf("%w, directory cannot be imported", ErrInvalidPath)
		}
		marker.TargetPath = input
	}

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
func processTargetDetail(marker *Marker, input string) error {
	exportMarker := regexp.MustCompile(`\[(\S+)\]`)

	markerRegex := exportMarker.FindStringSubmatch(input)
	switch {
	// Handle export marker
	case markerRegex != nil:
		marker.TargetExportMarker = string(markerRegex[1])

	// Handle line range marker with commas
	case strings.Contains(input, ","):
		targetLines := []int{}
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
		lb, ub, err := getLineRangeWithTilde(input)
		if err != nil {
			return fmt.Errorf("%w, %v", ErrInvalidSyntax, err)
		}
		marker.TargetLineFrom = lb
		marker.TargetLineTo = ub

	default:
		i, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("%w, %v", ErrInvalidSyntax, err)
		}
		marker.TargetLines = append(marker.TargetLines, i)
	}

	return nil
}

var (
	errLowerBound = errors.New("invalid lower bound in line range")
	errUpperBound = errors.New("invalid upper bound in line range")
)

func getLineRangeWithTilde(input string) (int, int, error) {
	lowerBound := 0
	upperBound := math.MaxInt32

	ls := strings.Split(input, "~")
	if len(ls) > 2 {
		return lowerBound, upperBound, fmt.Errorf("%w, tilde cannot be used more than once", ErrInvalidSyntax)
	}

	lb := ls[0]
	ub := ls[1]

	if lb != "" {
		l, err := strconv.Atoi(lb)
		if err != nil {
			return lowerBound, upperBound, fmt.Errorf("%w, %v", errUpperBound, err)
		}
		lowerBound = l
	}

	if ub != "" {
		u, err := strconv.Atoi(ub)
		if err != nil {
			return lowerBound, upperBound, fmt.Errorf("%w, %v", errLowerBound, err)
		}
		upperBound = u
	}

	return lowerBound, upperBound, nil
}

package parse

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// File holds onto file data.
type File struct {
	FileName string

	// fileType is derived from FileName, which is simply represented using
	// extension format.
	fileType string

	// contentBefore holds the file content as it was before processing. The
	// first slice represents the line number, and the second is for the actual
	// data.
	contentBefore [][]byte

	// contentPurged holds the file coontent, but removes the parts between
	// importer annotation begin/end. The first slice represents the line
	// number, and the second is for the actual data.
	contentPurged [][]byte

	// contentAfter holds the file content after the import has been run. This
	// only holds the actual data in byte slice representation.
	contentAfter []byte

	// annotations is an array holding onto each annotation block.
	annotations map[int]annotation
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
func Parse(fileName string, input io.Reader) (*File, error) {
	if input == nil {
		return nil, errors.New("no file content found")
	}

	fileType := filepath.Ext(fileName)

	fc := &File{
		FileName: fileName,
		fileType: fileType,
	}

	switch fileType {
	case ".md":
		if err := fc.parseMarkdown(input); err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		fmt.Printf("yaml")
	default:
		return nil, fmt.Errorf("unsupported file type '%s' provided", fc.fileType)
	}

	if err := fc.updateWithAnnotations(); err != nil {
		return nil, err
	}

	return fc, nil
}

type annotation struct {
	name                  string
	lineWithBeginOriginal int
	lineWithBeginPurged   int
	targetPath            string
	targetLines           []int
}

// matchHolder is a temporary data holder, which is used to ensure validity of
// annotation data.
type matchHolder struct {
	isBeginFound          bool
	isEndFound            bool
	lineWithBeginOriginal int
	lineWithBeginPurged   int
	option                []string
}

// parseMarkdown reads markdown input using scanner. This reads the input line
// by line, and store the data into File data. Parsing the data stores 3 sets
// of data: file content as is, annotation details, and file content with all
// data between annotation pairs purged.
func (f *File) parseMarkdown(input io.Reader) error {
	re := regexp.MustCompile(AnnotationMarkdown)

	annotations := map[int]annotation{}
	matches := map[string]matchHolder{}

	currentLine := 0
	inNested := false // Flag to check if the data is between annotations
	nestedUnder := "" // Name to check for annotation pair ending

	// NOTE:
	// For *File.contentXyz, I'm purposely making the first item in slice empty
	// for readability. This shouldn't be necessary, but with this approach,
	// the slice index matches the line number, and is easy to get my head
	// around for now.
	f.contentBefore = make([][]byte, 1)
	f.contentPurged = make([][]byte, 1)

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		currentLine++
		currentBytes := scanner.Bytes()
		f.contentBefore = append(f.contentBefore, currentBytes)

		match := re.FindStringSubmatch(scanner.Text())
		if len(match) == 0 {
			// If the line appears after some annotation, add empty slice.
			if inNested {
				continue
			}
			f.contentPurged = append(f.contentPurged, currentBytes)
			continue
		}

		// Parse regex match into groups to handle annotation
		subgroupName := match[1] // regex 1st subgroup. Index 0 is for full string.

		// Ensure this is the top most annotation. If an annotation is found
		// within another annotation as nested, ignore it. This is because we
		// should be handling those nested annotations in those target files
		// instead.
		// This means, in any file in question, the parse logic only looks at
		// one file and its direct dependencies.
		if inNested && nestedUnder != subgroupName {
			f.contentPurged = append(f.contentPurged, nil) // Add empty slice
			continue
		}

		nestedUnder = subgroupName

		// At this point, the annotation is important, and we need to process
		// the line further.
		// Note that, contentPurged does not contain any data that's wrapped
		// between annotations. Those lines will be kept as an empty byte slice
		// for further processing later to create contentAfter.
		f.contentPurged = append(f.contentPurged, currentBytes)

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
					matchData.lineWithBeginOriginal = currentLine
					matchData.lineWithBeginPurged = len(f.contentPurged)
				}
				if matchedContent == "end" {
					inNested = false
					nestedUnder = ""
					matchData.isEndFound = true
				}
			case "importer_option":
				matchData.option = append(matchData.option, string(matchedContent))
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

		annotations[annotation.lineWithBeginPurged] = annotation
	}

	f.annotations = annotations

	return nil
}

func convert(name string, match matchHolder) (annotation, error) {
	if !match.isBeginFound || !match.isEndFound {
		return annotation{}, errors.New("no matching annotations found")
	}
	reImportTarget := regexp.MustCompile(OptionFilePathIndicator)

	result := annotation{
		name:                  name,
		lineWithBeginOriginal: match.lineWithBeginOriginal,
		lineWithBeginPurged:   match.lineWithBeginPurged,
	}

	for _, opt := range match.option {
		match := reImportTarget.FindAllStringSubmatch(opt, -1)
		if len(match) == 0 {
			continue
		}

		for _, ms := range match {
			for i, n := range reImportTarget.SubexpNames() {
				matchedContent := ms[i]
				switch n {
				case "importer_target_path":
					result.targetPath = matchedContent
				case "importer_target_lines":
					lines := []int{}

					// Handle case where the input is something like "6~22"
					if strings.Contains(matchedContent, "~") {
						ls := strings.Split(matchedContent, "~")
						lb, err := strconv.Atoi(ls[0])
						if err != nil {
							fmt.Printf("error: found non-numeric input for line number lower bound, %v", err)
						}
						ub, err := strconv.Atoi(ls[1])
						if err != nil {
							fmt.Printf("error: found non-numeric input for line number upper bound, %v", err)
						}
						// Add line numbers to the slice.
						// This way, we can support comma separated list, etc.
						for i := lb; i <= ub; i++ {
							lines = append(lines, i)
						}
					}
					result.targetLines = lines
				}
			}
		}
	}

	return result, nil
}

func (f *File) updateWithAnnotations() error {
	result := []byte{}
	br := byte('\n')
	for line, data := range f.contentPurged {
		result = append(result, data...)
		result = append(result, br)

		if a, found := f.annotations[line+1]; found {
			f, err := os.Open(a.targetPath)
			if err != nil {
				fmt.Printf("error: could not open file '%s'", a.targetPath)
				continue
			}
			currentLine := 0
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				currentLine++
				for _, l := range a.targetLines {
					if currentLine == l {
						result = append(result, scanner.Bytes()...)
						result = append(result, br)
					}
				}
			}
		}
	}
	f.contentAfter = result
	return nil
}

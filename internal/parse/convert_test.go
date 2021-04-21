package parse

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/file"
)

func TestConvert(t *testing.T) {
	cases := map[string]struct {
		// Input
		name  string
		match matchHolder

		// Output
		wantResult *file.Annotation
		wantErr    error
	}{
		"valid annotation match": {
			name: "test name",
			match: matchHolder{
				isBeginFound:   true,
				isEndFound:     true,
				lineToInsertAt: 10,
				options:        "from: ./some_file.txt#2~22",
			},
			wantResult: &file.Annotation{
				Name:           "test name",
				LineToInsertAt: 10,
				TargetPath:     "./some_file.txt",
				TargetLineFrom: 2,
				TargetLineTo:   22,
			},
		},
		"TBC: annotation without option, valid for now": {
			name: "test name",
			match: matchHolder{
				isBeginFound:   true,
				isEndFound:     true,
				lineToInsertAt: 10,
			},
			wantResult: &file.Annotation{
				Name:           "test name",
				LineToInsertAt: 10,
			},
		},

		// ERROR CASES
		"INVALID: annotation is not matched, end missing": {
			name: "test name",
			match: matchHolder{
				isBeginFound:   true,
				isEndFound:     false,
				lineToInsertAt: 10,
				options:        "from: ./some_file.txt#2~22",
			},
			wantErr: ErrNoMatchingAnnotations,
		},
		"INVALID: annotation is not matched, beging missing": {
			name: "test name",
			match: matchHolder{
				isBeginFound:   false,
				isEndFound:     true,
				lineToInsertAt: 10,
				options:        "from: ./some_file.txt#2~22",
			},
			wantErr: ErrNoMatchingAnnotations,
		},
		"INVALID: annotation has invalid line ranges": {
			name: "test name",
			match: matchHolder{
				isBeginFound:   true,
				isEndFound:     true,
				lineToInsertAt: 10,
				options:        "from: file.txt#x~y",
			},
			wantErr: ErrInvalidSyntax,
		},
		// Commenting out, as there is no code path that generates invalid path, as it would usually be ignored by regex
		// "INVALID: annotation has invalid path": {
		// 	name: "test name",
		// 	match: matchHolder{
		// 		isBeginFound:   true,
		// 		isEndFound:     true,
		// 		lineToInsertAt: 10,
		// 		options:        "from: some_invalid_path#2~22", // There is no validation as of now
		// 	},
		// 	wantErr: ErrNoMatchingAnnotations,
		// },
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			annotation, err := convert(tc.name, tc.match)
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}

			if diff := cmp.Diff(tc.wantResult, annotation, cmp.AllowUnexported(file.Annotation{})); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestProcessTargetPath(t *testing.T) {
	cases := map[string]struct {
		// Input
		annotation *file.Annotation
		input      string

		// Output
		wantResult *file.Annotation
		wantErr    error
	}{
		"simple test": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input: "./some_path.txt",
			wantResult: &file.Annotation{
				Name:       "test file",
				TargetPath: "./some_path.txt",
			},
		},

		"INVALID: empty file path": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input:   "",
			wantErr: ErrInvalidPath,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := processTargetPath(tc.annotation, tc.input)
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}

			if diff := cmp.Diff(tc.wantResult, tc.annotation, cmp.AllowUnexported(file.Annotation{})); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestProcessTargetDetail(t *testing.T) {
	cases := map[string]struct {
		// Input
		annotation *file.Annotation
		input      string

		// Output
		wantResult *file.Annotation
		wantErr    error
	}{
		"export marker": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input: "[some_marker]",
			wantResult: &file.Annotation{
				Name:               "test file",
				TargetExportMarker: "some_marker",
			},
		},
		"simple comma separated values": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input: "1,2,3,5",
			wantResult: &file.Annotation{
				Name:        "test file",
				TargetLines: []int{1, 2, 3, 5},
			},
		},
		"simple tilde based values": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input: "3~22",
			wantResult: &file.Annotation{
				Name:           "test file",
				TargetLineFrom: 3,
				TargetLineTo:   22,
			},
		},
		"tilde based values without lower bound": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input: "~22",
			wantResult: &file.Annotation{
				Name:         "test file",
				TargetLineTo: 22,
			},
		},
		"tilde based values without upper bound": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input: "3~",
			wantResult: &file.Annotation{
				Name:           "test file",
				TargetLineFrom: 3,
			},
		},
		"comma separated, complex values": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input: "2,1,1,1,3~7,4~6",
			wantResult: &file.Annotation{
				Name:        "test file",
				TargetLines: []int{2, 1, 1, 1, 3, 4, 5, 6, 7, 4, 5, 6},
			},
		},
		"single line value": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input: "15",
			wantResult: &file.Annotation{
				Name:        "test file",
				TargetLines: []int{15},
			},
		},

		// ERROR CASE
		"INVALID: export marker": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input:   "[some marker with whitespace]",
			wantErr: ErrInvalidSyntax,
		},
		"INVALID: tilde based range, invalid char lower range": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input:   "x~2",
			wantErr: ErrInvalidSyntax,
		},
		"INVALID: tilde based range, invalid char upper range": {
			annotation: &file.Annotation{
				Name: "test file",
			},
			input:   "3~x",
			wantErr: ErrInvalidSyntax,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := processTargetDetail(tc.annotation, tc.input)
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}

			if diff := cmp.Diff(tc.wantResult, tc.annotation, cmp.AllowUnexported(file.Annotation{})); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}
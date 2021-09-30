package parse

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/file"
	"github.com/upsidr/importer/internal/marker"
	"github.com/upsidr/importer/internal/testingutil/golden"
)

func TestParseMarkdown(t *testing.T) {
	cases := map[string]struct {
		// Input
		fileName string
		input    io.Reader

		// Output
		wantFile *file.File
	}{
		"simple test from main testdata": {
			fileName: "dummy.md",
			input:    strings.NewReader(golden.FileAsString(t, "../../testdata/markdown/simple-before.md")),
			wantFile: &file.File{
				FileName: "dummy.md",
				ContentBefore: StringToLineStrings(t,
					golden.FileAsString(t, "../../testdata/markdown/simple-before.md")),
				ContentPurged: StringToLineStrings(t,
					golden.FileAsString(t, "../../testdata/markdown/simple-purged.md")),
				Markers: map[int]*marker.Marker{
					3: {
						Name:           "lorem",
						LineToInsertAt: 3,
						ImportTargetFile: marker.ImportTargetFile{
							Type: marker.PathBased,
							File: "./snippet-lorem.md",
						},
						ImportLogic: marker.ImportLogic{
							Type:     marker.LineRange,
							LineFrom: 5,
							LineTo:   12,
						},
					},
				},
			},
		},
		"no importer annotation": {
			fileName: "./testdata/markdown/no-importer-marker-before.md",
			wantFile: &file.File{
				FileName: "./testdata/markdown/no-importer-marker-before.md",
				ContentBefore: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/markdown/no-importer-marker-before.md")),
				ContentPurged: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/markdown/no-importer-marker-purged.md")),
				Markers: map[int]*marker.Marker{},
			},
		},
		"with single importer annotation": {
			fileName: "./testdata/markdown/single-marker-before.md",
			wantFile: &file.File{
				FileName: "./testdata/markdown/single-marker-before.md",
				ContentBefore: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/markdown/single-marker-before.md")),
				ContentPurged: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/markdown/single-marker-purged.md")),
				Markers: map[int]*marker.Marker{
					3: {
						Name:           "some_importer",
						LineToInsertAt: 3,
						ImportTargetFile: marker.ImportTargetFile{
							Type: marker.PathBased,
							File: "../../testdata/markdown/simple-before-importer.md",
						},
						ImportLogic: marker.ImportLogic{
							Type:     marker.LineRange,
							LineFrom: 1,
							LineTo:   2,
						},
					},
				},
			},
		},
		"with single importer annotation, inner annotation ignored": {
			fileName: "./testdata/markdown/single-marker-with-inner-before.md",
			wantFile: &file.File{
				FileName: "./testdata/markdown/single-marker-with-inner-before.md",
				ContentBefore: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/markdown/single-marker-with-inner-before.md")),
				ContentPurged: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/markdown/single-marker-with-inner-purged.md")),
				Markers: map[int]*marker.Marker{
					3: {
						Name:           "some_importer",
						LineToInsertAt: 3,
						ImportTargetFile: marker.ImportTargetFile{
							Type: marker.PathBased,
							File: "./somefile",
						},
						ImportLogic: marker.ImportLogic{
							Type:     marker.LineRange,
							LineFrom: 1,
							LineTo:   2,
						},
					},
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fileInput := tc.input
			if fileInput == nil {
				fileInput = golden.FileAsReader(t, tc.fileName)
			}
			f, err := Parse(tc.fileName, fileInput)
			if err != nil {
				t.Errorf("unexpected error, %v", err)
				return
			}

			if diff := cmp.Diff(tc.wantFile, f, cmp.AllowUnexported(file.File{})); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestParseMarkdownFail(t *testing.T) {
	cases := map[string]struct {
		// Input
		fileName string
		input    io.Reader

		// Output
		wantErr error
	}{
		// ===================
		// Invalid cases below
		"importer option missing": {
			fileName: "dummy.md",
			input: strings.NewReader(`
# Test Markdown

<!-- == imptr: some_importer / begin == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
			wantErr: marker.ErrInvalidSyntax,
		},
		"file line range not number - lower bound": {
			fileName: "dummy.md",
			input: strings.NewReader(`
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#NOT_NUMBER~2233 == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
			wantErr: marker.ErrInvalidSyntax,
		},
		"file line range not number - upper bound": {
			fileName: "dummy.md",
			input: strings.NewReader(`
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#1~NOT_NUMBER == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
			wantErr: marker.ErrInvalidSyntax,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fileInput := tc.input
			if fileInput == nil {
				fileInput = golden.FileAsReader(t, tc.fileName)
			}
			_, err := Parse(tc.fileName, fileInput)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
			}
		})
	}
}

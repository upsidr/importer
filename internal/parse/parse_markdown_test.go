package parse

import (
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/file"
)

func TestParseMarkdown(t *testing.T) {
	cases := map[string]struct {
		// Input
		fileName string
		input    io.Reader

		// Output
		wantFile *file.File
	}{
		"no importer annotation": {
			fileName: "dummy",
			input: strings.NewReader(`
# Test Markdown

No importer annotation
`),
			wantFile: &file.File{
				FileName: "dummy",
				ContentBefore: StringToLineBytes(t, `
# Test Markdown

No importer annotation
`),
				ContentPurged: StringToLineBytes(t, `
# Test Markdown

No importer annotation
`),
				Annotations: map[int]*file.Annotation{},
			},
		},
		"with single importer annotation": {
			fileName: "dummy",
			input: strings.NewReader(`
# Test Markdown

<!-- == imptr: some_importer / begin from: ../../testdata/simple-before-importer.md#1~2 == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
			wantFile: &file.File{
				FileName: "dummy",
				ContentBefore: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ../../testdata/simple-before-importer.md#1~2 == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
				ContentPurged: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ../../testdata/simple-before-importer.md#1~2 == -->
<!-- == imptr: some_importer / end == -->
`),
				Annotations: map[int]*file.Annotation{
					4: {
						Name:           "some_importer",
						LineToInsertAt: 4,
						TargetPath:     "../../testdata/simple-before-importer.md",
						TargetLines:    []int{1, 2},
					},
				},
			},
		},
		"with single importer annotation, inner annotation ignored": {
			fileName: "dummy",
			input: strings.NewReader(`
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#1~2 == -->

some data between an annotation pair, which gets purged.

This annotation for "another_importer" gets ignored as it is within another annotation pair.
<!-- == imptr: another_importer / begin from: ./another_file#1~2 == -->
<!-- == imptr: another_importer / end == -->

<!-- == imptr: some_importer / end == -->
`),
			wantFile: &file.File{
				FileName: "dummy",
				ContentBefore: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#1~2 == -->

some data between an annotation pair, which gets purged.

This annotation for "another_importer" gets ignored as it is within another annotation pair.
<!-- == imptr: another_importer / begin from: ./another_file#1~2 == -->
<!-- == imptr: another_importer / end == -->

<!-- == imptr: some_importer / end == -->
`),
				ContentPurged: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#1~2 == -->
<!-- == imptr: some_importer / end == -->
`),
				Annotations: map[int]*file.Annotation{
					4: {
						Name:           "some_importer",
						LineToInsertAt: 4,
						TargetPath:     "./somefile",
						TargetLines:    []int{1, 2},
					},
				},
			},
		},

		// ===================
		// Invalid cases below
		"importer option missing": {
			fileName: "dummy",
			input: strings.NewReader(`
# Test Markdown

<!-- == imptr: some_importer / begin == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
			wantFile: &file.File{
				FileName: "dummy",
				ContentBefore: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
				ContentPurged: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin == -->
<!-- == imptr: some_importer / end == -->
`),
				Annotations: map[int]*file.Annotation{
					4: {
						// Name of improter and line are found
						Name:           "some_importer",
						LineToInsertAt: 4,
						// But no target specified
					},
				},
			},
		},
		"file line range not number - lower bound": {
			fileName: "dummy",
			input: strings.NewReader(`
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#NOT_NUMBER~2233 == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
			wantFile: &file.File{
				FileName: "dummy",
				ContentBefore: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#NOT_NUMBER~2233 == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
				ContentPurged: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#NOT_NUMBER~2233 == -->
<!-- == imptr: some_importer / end == -->
`),
				Annotations: map[int]*file.Annotation{},
			},
		},
		"file line range not number - upper bound": {
			fileName: "dummy",
			input: strings.NewReader(`
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#1~NOT_NUMBER == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
			wantFile: &file.File{
				FileName: "dummy",
				ContentBefore: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#1~NOT_NUMBER == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
				ContentPurged: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ./somefile#1~NOT_NUMBER == -->
<!-- == imptr: some_importer / end == -->
`),
				Annotations: map[int]*file.Annotation{},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f, err := parseMarkdown(tc.fileName, tc.input)
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

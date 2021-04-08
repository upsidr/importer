package parse

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/upsidr/importer/internal/file"
)

func TestParse(t *testing.T) {
	cases := map[string]struct {
		// Input
		fileName string
		input    io.Reader

		// Output
		wantFile *file.File
		wantErr  error
	}{
		"markdown - more extensive tests in parse_markdown_test.go": {
			fileName: "dummy.md",
			input: strings.NewReader(`
# Test Markdown

<!-- == imptr: some_importer / begin from: ../../testdata/note.txt#1~3 == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
			wantFile: &file.File{
				FileName: "dummy.md",
				ContentBefore: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ../../testdata/note.txt#1~3 == -->
some data between an annotation pair, which gets purged.
<!-- == imptr: some_importer / end == -->
`),
				ContentPurged: StringToLineBytes(t, `
# Test Markdown

<!-- == imptr: some_importer / begin from: ../../testdata/note.txt#1~3 == -->
<!-- == imptr: some_importer / end == -->
`),
				ContentAfter: []byte(`
# Test Markdown

<!-- == imptr: some_importer / begin from: ../../testdata/note.txt#1~3 == -->
This is test data.
他言語サポートのためのテスト文章。
🍸 Emojis 🍷 Supported 🍺
<!-- == imptr: some_importer / end == -->
`),
				Annotations: map[int]*file.Annotation{
					4: {
						Name:           "some_importer",
						LineToInsertAt: 4,
						TargetPath:     "../../testdata/note.txt",
						TargetLines:    []int{1, 2, 3},
					},
				},
			},
		},

		// ERROR CASES
		"extension not supported": {
			fileName: "no_extension",
			input:    strings.NewReader("dummy"),
			wantErr:  ErrUnsupportedFileType,
		},
		"no data": {
			fileName: "dummy.md",
			input:    nil,
			wantErr:  ErrNoInput,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			f, err := Parse(tc.fileName, tc.input)
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}

			if diff := cmp.Diff(tc.wantFile, f, cmp.AllowUnexported(file.File{})); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

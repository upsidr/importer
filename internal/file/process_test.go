package file

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessAnnotation(t *testing.T) {
	cases := map[string]struct {
		// Input
		file *File

		// Output
		want    []byte
		wantErr error
	}{
		"test": {
			file: &File{
				FileName: "test-file.md",
				ContentPurged: []string{
					"This is",
					"a test",
					"data",
				},
				Annotations: map[int]*Annotation{
					2: {
						Name:           "test annotation",
						LineToInsertAt: 2,
						TargetPath:     "../../testdata/note.txt",
						TargetLines:    []int{1, 2},
					},
				},
			},
			want: []byte(`This is
a test
This is test data.
他言語サポートのためのテスト文章。
data
`),
		},
		"file does not exist, and gets ignored with a warning": {
			file: &File{
				FileName: "test-file.md",
				ContentPurged: []string{
					"This is",
					"a test",
					"data",
				},
				Annotations: map[int]*Annotation{
					2: {
						Name:           "test annotation",
						LineToInsertAt: 2,
						TargetPath:     "../../does-not-exist.txt",
						TargetLines:    []int{1, 2},
					},
				},
			},
			want: []byte(`This is
a test
data
`),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.file.ProcessAnnotations()
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}

			if diff := cmp.Diff(tc.want, tc.file.ContentAfter); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestProcessSingleAnnotation(t *testing.T) {
	cases := map[string]struct {
		// Input
		previousData []byte
		callerFile   string
		annotation   *Annotation

		// Output
		want    []byte
		wantErr error
	}{
		"range process": {
			previousData: []byte(`Some data
and another line
`),
			callerFile: "./some_file.md",
			annotation: &Annotation{
				LineToInsertAt: 1, // Not used in this, as single annotation handling is about appending data
				TargetPath:     "../../testdata/note.txt",
				TargetLineFrom: 1,
				TargetLineTo:   3,
			},
			want: []byte(`Some data
and another line
This is test data.
他言語サポートのためのテスト文章。
🍸 Emojis 🍷 Supported 🍺
`),
		},
		"comma separated lines": {
			previousData: []byte{},
			callerFile:   "./some_file.md",
			annotation: &Annotation{
				LineToInsertAt: 1,
				TargetPath:     "../../testdata/note.txt",
				TargetLines:    []int{2, 3},
			},
			want: []byte(`他言語サポートのためのテスト文章。
🍸 Emojis 🍷 Supported 🍺
`),
		},
		"exporter marker": {
			previousData: []byte{},
			callerFile:   "./some_file.md",
			annotation: &Annotation{
				LineToInsertAt:     1,
				TargetPath:         "../../testdata/with-exporter.md",
				TargetExportMarker: "test_exporter",
			},
			want: []byte(`
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨
🚀🚀🚀🚀🚀🚀🚀🚀
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨

`),
		},

		// ERROR CASES
		"no target file found": {
			previousData: []byte{},
			callerFile:   "./some_file.md",
			annotation: &Annotation{
				LineToInsertAt:     1,
				TargetPath:         "../../does-not-exist.md",
				TargetExportMarker: "test_exporter",
			},
			wantErr: os.ErrNotExist,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := processSingleAnnotation(tc.previousData, tc.callerFile, tc.annotation)
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}

			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

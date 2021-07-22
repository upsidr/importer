package file

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProcessMarker(t *testing.T) {
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
				Markers: map[int]*Marker{
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
				Markers: map[int]*Marker{
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
			err := tc.file.ProcessMarkers()
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

func TestProcessSingleMarker(t *testing.T) {
	cases := map[string]struct {
		// Input
		callerFile string
		marker     *Marker

		// Output
		want    []byte
		wantErr error
	}{
		"markdown: range process": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt: 1, // Not used in this, as single annotation handling is about appending data
				TargetPath:     "../../testdata/note.txt",
				TargetLineFrom: 1,
				TargetLineTo:   3,
			},
			want: []byte(`This is test data.
他言語サポートのためのテスト文章。
🍸 Emojis 🍷 Supported 🍺
`),
		},
		"markdown: comma separated lines": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt: 1,
				TargetPath:     "../../testdata/note.txt",
				TargetLines:    []int{2, 3},
			},
			want: []byte(`他言語サポートのためのテスト文章。
🍸 Emojis 🍷 Supported 🍺
`),
		},
		"markdown: exporter marker": {
			callerFile: "./some_file.md",
			marker: &Marker{
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
		"yaml: simple import": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt:     5,
				TargetPath:         "../../testdata/with-exporter.yaml",
				TargetExportMarker: "long-tree",
				Indentation:        nil,
			},
			want: []byte(`a:
  b:
    c:
      d:
        e:
          f:
            g:
              h:
                i:
                  j:
                    k: {}
`),
		},
		"yaml: absolute indentation": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt:     5,
				TargetPath:         "../../testdata/with-exporter.yaml",
				TargetExportMarker: "metadata-only",
				Indentation: &Indentation{
					Mode:   AbsoluteIndentation,
					Length: 30,
				},
			},
			want: []byte(`                              metadata:
                                name: sample-data
                                namespace: sample-namespace
`),
		},
		"yaml: absolute indentation with zero indentation": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt:     5,
				TargetPath:         "../../testdata/with-exporter.yaml",
				TargetExportMarker: "metadata-only",
				Indentation: &Indentation{
					Mode:   AbsoluteIndentation,
					Length: 0,
				},
			},
			want: []byte(`metadata:
  name: sample-data
  namespace: sample-namespace
`),
		},
		"yaml: extra indentation": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt:     5,
				TargetPath:         "../../testdata/with-exporter.yaml",
				TargetExportMarker: "sample-nested",
				Indentation: &Indentation{
					Mode:   ExtraIndentation,
					Length: 2,
				},
			},
			want: []byte(`    nested:
      more:
        data:
          sample: This is a sample data
        metadata:
          name: sample-data
          namespace: sample-namespace
`),
		},

		// ERROR CASES
		"no target file found": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt:     1,
				TargetPath:         "../../does-not-exist.md",
				TargetExportMarker: "test_exporter",
			},
			wantErr: os.ErrNotExist,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := processSingleMarker(tc.callerFile, tc.marker)
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

func TestPrependWhitespaces(t *testing.T) {
	cases := map[string]struct {
		originalSlice   []byte
		whitespaceCount int

		want []byte
	}{
		"indentation": {
			originalSlice:   []byte("abcdef"),
			whitespaceCount: 6,
			want:            []byte("      abcdef"),
		},
		"extra indentation": {
			originalSlice:   []byte("  abcdef"),
			whitespaceCount: 6,
			want:            []byte("        abcdef"),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := prependWhitespaces(tc.originalSlice, tc.whitespaceCount)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("prepend result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

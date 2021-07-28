package marker

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

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
				TargetPath:     "../../testdata/other/note.txt",
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
				TargetPath:     "../../testdata/other/note.txt",
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
				TargetPath:         "../../testdata/markdown/with-exporter.md",
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
		"yaml: line range": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				TargetPath:     "../../testdata/yaml/simple-tree.yaml",
				TargetLineFrom: 2,
				TargetLineTo:   5,
				Indentation:    nil,
			},
			want: []byte(`  b:
    c:
      d:
        e:
`),
		},
		"yaml: comma separated lines": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				TargetPath:     "../../testdata/yaml/simple-tree.yaml",
				TargetLines:    []int{1, 2, 4},
				Indentation:    nil,
			},
			want: []byte(`a:
  b:
      d:
`),
		},
		"yaml: exporter marker": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt:     5,
				TargetPath:         "../../testdata/yaml/with-exporter.yaml",
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
		"yaml: exporter marker with absolute indentation": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt:     5,
				TargetPath:         "../../testdata/yaml/with-exporter.yaml",
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
		"yaml: exporter marker with absolute indentation, with zero indentation": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt:     5,
				TargetPath:         "../../testdata/yaml/with-exporter.yaml",
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
		"yaml: exporter marker with extra indentation": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt:     5,
				TargetPath:         "../../testdata/yaml/with-exporter.yaml",
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
		"other: range process": {
			callerFile: "./some_unknown_file_type",
			marker: &Marker{
				LineToInsertAt: 1, // Not used in this, as single annotation handling is about appending data
				TargetPath:     "../../testdata/other/note.txt",
				TargetLineFrom: 1,
				TargetLineTo:   3,
			},
			want: []byte(`This is test data.
他言語サポートのためのテスト文章。
🍸 Emojis 🍷 Supported 🍺
`),
		},
		"other: comma separated lines": {
			callerFile: "./some_unknown_file_type",
			marker: &Marker{
				LineToInsertAt: 1,
				TargetPath:     "../../testdata/other/note.txt",
				TargetLines:    []int{2, 3},
			},
			want: []byte(`他言語サポートのためのテスト文章。
🍸 Emojis 🍷 Supported 🍺
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
			result, err := tc.marker.ProcessMarker(tc.callerFile)
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

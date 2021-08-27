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
				TargetPath:         "../../testdata/markdown/snippet-with-exporter.md",
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
				TargetPath:     "../../testdata/yaml/snippet-simple-tree.yaml",
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
		"yaml: line range with URL, github.com": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				TargetURL:      "https://github.com/upsidr/importer/blob/main/testdata/yaml/snippet-simple-tree.yaml",
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
		"yaml: line range with URL, raw.githubusercontent.com": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				TargetURL:      "https://raw.githubusercontent.com/upsidr/importer/main/testdata/yaml/snippet-simple-tree.yaml",
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
				TargetPath:     "../../testdata/yaml/snippet-simple-tree.yaml",
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
				TargetPath:         "../../testdata/yaml/snippet-with-exporter.yaml",
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
				TargetPath:         "../../testdata/yaml/snippet-with-exporter.yaml",
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
				TargetPath:         "../../testdata/yaml/snippet-with-exporter.yaml",
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
				TargetPath:         "../../testdata/yaml/snippet-with-exporter.yaml",
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
		"yaml: exporter marker with align indentation": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt:     5,
				TargetPath:         "../../testdata/yaml/snippet-with-exporter.yaml",
				TargetExportMarker: "sample-nested",
				Indentation: &Indentation{
					Mode:              AlignIndentation,
					MarkerIndentation: 10,
				},
			},
			want: []byte(`          nested:
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
		"no target file information provided": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt:     1,
				TargetPath:         "", // important
				TargetURL:          "", // important
				TargetExportMarker: "test_exporter",
			},
			wantErr: ErrNoFileInput,
		},
		"url access error": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				TargetURL:      "https://some-address-that-does-not-exist",
				TargetLineFrom: 1,
				TargetLineTo:   5,
				Indentation:    nil,
			},
			wantErr: ErrGetMarkerTarget,
		},
		"404 not found error with github.com": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				TargetURL:      "https://github.com/does-not-exist.yml#",
				TargetLineFrom: 1,
				TargetLineTo:   5,
				Indentation:    nil,
			},
			wantErr: ErrNonSuccessCode,
		},
		"invalid address": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				TargetURL:      "http///////",
				TargetLineFrom: 1,
				TargetLineTo:   5,
				Indentation:    nil,
			},
			wantErr: ErrInvalidURL,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := tc.marker.ProcessMarkerData(tc.callerFile)
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}

			if diff := cmp.Diff(string(tc.want), string(result)); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestPreprocessURL(t *testing.T) {
	cases := map[string]struct {
		input   string
		want    string
		wantErr error
	}{
		"github.com -> raw: based on branch": {
			input: "https://github.com/upsidr/importer/blob/main/README.md",
			want:  "https://raw.githubusercontent.com/upsidr/importer/main/README.md",
		},
		"github.com -> raw: based on commit": {
			input: "https://github.com/upsidr/importer/blob/a81365f65633be66972c5b071ba9db6bc44ffeb3/README.md",
			want:  "https://raw.githubusercontent.com/upsidr/importer/a81365f65633be66972c5b071ba9db6bc44ffeb3/README.md",
		},
		"github.com -> raw: directory 'blob' is ignored": {
			input: "https://github.com/upsidr/importer/blob/main/blob/blob/README.md",
			want:  "https://raw.githubusercontent.com/upsidr/importer/main/blob/blob/README.md",
		},
		"github.com as is: directory input": {
			input: "https://github.com/upsidr/importer/tree/main/cmd/importer",
			want:  "https://github.com/upsidr/importer/tree/main/cmd/importer",
		},
		"github.com as is: some other content": {
			input: "https://github.com/explore",
			want:  "https://github.com/explore",
		},
		"non-github.com: return as is": {
			input: "https://google.com/something",
			want:  "https://google.com/something",
		},

		// ERROR CASE
		"invalid path": {
			input:   "something",
			wantErr: ErrInvalidPath,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := preprocessURL(tc.input)
			if (err != nil) && !errors.Is(err, tc.wantErr) {
				t.Fatalf("error mismatch\n    error = %v\n    wantErr %v", err, tc.wantErr)
			}

			if diff := cmp.Diff(string(tc.want), string(result)); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

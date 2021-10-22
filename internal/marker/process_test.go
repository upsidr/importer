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
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/other/note.txt",
				},
				ImportLogic: ImportLogic{
					Type:     LineRange,
					LineFrom: 1,
					LineTo:   3,
				},
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
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/other/note.txt",
				},
				ImportLogic: ImportLogic{
					Type:  CommaSeparatedLines,
					Lines: []int{2, 3},
				},
			},
			want: []byte(`他言語サポートのためのテスト文章。
🍸 Emojis 🍷 Supported 🍺
`),
		},
		"markdown: exporter marker": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt: 1,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/markdown/snippet-with-exporter.md",
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "test_exporter",
				},
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
		"markdown: exporter marker with quote": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt: 1,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/markdown/snippet-with-exporter.md",
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "test_exporter",
				},
				ImportStyle: &ImportStyle{
					Mode: Quote,
				},
			},
			want: []byte(`> 
> ✨✨✨✨✨✨✨✨
> ✨✨✨✨✨✨✨✨
> ✨✨✨✨✨✨✨✨
> 🚀🚀🚀🚀🚀🚀🚀🚀
> ✨✨✨✨✨✨✨✨
> ✨✨✨✨✨✨✨✨
> ✨✨✨✨✨✨✨✨
> 
`),
		},
		"markdown: exporter marker with verbatim": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt: 1,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/markdown/snippet-with-exporter.md",
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "test_exporter",
				},
				Wrap: &Wrap{
					LanguageType: "some-lang",
				},
			},
			want: []byte("```" + `some-lang

✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨
🚀🚀🚀🚀🚀🚀🚀🚀
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨
✨✨✨✨✨✨✨✨

` + "```" + `
`),
		},
		"yaml: line range": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/yaml/snippet-simple-tree.yaml",
				},
				ImportLogic: ImportLogic{
					Type:     LineRange,
					LineFrom: 2,
					LineTo:   5,
				},
				Indentation: nil,
			},
			want: []byte(`  b:
    c:
      d:
        e:
`),
		},
		"yaml: line range with URL, github.com - this may fail if GitHub is down": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: URLBased,
					File: "https://github.com/upsidr/importer/blob/main/testdata/yaml/snippet-simple-tree.yaml",
				},
				ImportLogic: ImportLogic{
					Type:     LineRange,
					LineFrom: 2,
					LineTo:   5,
				},
				Indentation: nil,
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
				ImportTargetFile: ImportTargetFile{
					Type: URLBased,
					File: "https://raw.githubusercontent.com/upsidr/importer/main/testdata/yaml/snippet-simple-tree.yaml",
				},
				ImportLogic: ImportLogic{
					Type:     LineRange,
					LineFrom: 2,
					LineTo:   5,
				},
				Indentation: nil,
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
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/yaml/snippet-simple-tree.yaml",
				},
				ImportLogic: ImportLogic{
					Type:  LineRange,
					Lines: []int{1, 2, 4},
				},
				Indentation: nil,
			},
			want: []byte(`a:
  b:
      d:
`),
		},
		"yaml: exporter marker": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/yaml/snippet-with-exporter.yaml",
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "long-tree",
				},
				Indentation: nil,
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
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/yaml/snippet-with-exporter.yaml",
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "metadata-only",
				},
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
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/yaml/snippet-with-exporter.yaml",
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "metadata-only",
				},
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
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/yaml/snippet-with-exporter.yaml",
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "sample-nested",
				},
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
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/yaml/snippet-with-exporter.yaml",
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "sample-nested",
				},
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
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/other/note.txt",
				},
				ImportLogic: ImportLogic{
					Type:     LineRange,
					LineFrom: 1,
					LineTo:   3,
				},
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
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../testdata/other/note.txt",
				},
				ImportLogic: ImportLogic{
					Type:  LineRange,
					Lines: []int{2, 3},
				},
			},
			want: []byte(`他言語サポートのためのテスト文章。
🍸 Emojis 🍷 Supported 🍺
`),
		},

		// ERROR CASES
		"no target file found": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt: 1,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "../../does-not-exist.md",
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "test_exporter",
				},
			},
			wantErr: os.ErrNotExist,
		},
		"no target file information provided - path": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt: 1,
				ImportTargetFile: ImportTargetFile{
					Type: PathBased,
					File: "", // Important
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "test_exporter",
				},
			},
			wantErr: ErrNoFileInput,
		},
		"no target file information provided - URL": {
			callerFile: "./some_file.md",
			marker: &Marker{
				LineToInsertAt: 1,
				ImportTargetFile: ImportTargetFile{
					Type: URLBased,
					File: "", // Important
				},
				ImportLogic: ImportLogic{
					Type:           ExporterMarker,
					ExporterMarker: "test_exporter",
				},
			},
			wantErr: ErrInvalidURL,
		},
		"url access error": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: URLBased,
					File: "https://some-address-that-does-not-exist",
				},
				ImportLogic: ImportLogic{
					Type:     LineRange,
					LineFrom: 1,
					LineTo:   5,
				},
				Indentation: nil,
			},
			wantErr: ErrGetMarkerTarget,
		},
		"404 not found error with github.com": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: URLBased,
					File: "https://github.com/does-not-exist.yml#",
				},
				ImportLogic: ImportLogic{
					Type:     LineRange,
					LineFrom: 1,
					LineTo:   5,
				},
				Indentation: nil,
			},
			wantErr: ErrNonSuccessCode,
		},
		"invalid address": {
			callerFile: "./some_file.yaml",
			marker: &Marker{
				LineToInsertAt: 5,
				ImportTargetFile: ImportTargetFile{
					Type: URLBased,
					File: "http///////",
				},
				ImportLogic: ImportLogic{
					Type:     LineRange,
					LineFrom: 1,
					LineTo:   5,
				},
				Indentation: nil,
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

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

func TestParseYAML(t *testing.T) {
	cases := map[string]struct {
		// Input
		fileName string
		input    io.Reader

		// Output
		wantFile *file.File
	}{
		"single importer annotation": {
			fileName: "./testdata/yaml/single-marker-before.yaml",
			wantFile: &file.File{
				FileName: "./testdata/yaml/single-marker-before.yaml",
				ContentBefore: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/yaml/single-marker-before.yaml")),
				ContentPurged: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/yaml/single-marker-purged.yaml")),
				Markers: map[int]*marker.Marker{
					3: {
						Name:               "some-importer",
						LineToInsertAt:     3,
						TargetPath:         "./exporter-example.yaml",
						TargetExportMarker: "random-data",
					},
				},
			},
		},
		"no importer annotation": {
			fileName: "./testdata/yaml/no-importer-marker-before.yaml",
			wantFile: &file.File{
				FileName: "./testdata/yaml/no-importer-marker-before.yaml",
				ContentBefore: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/yaml/no-importer-marker-before.yaml")),
				ContentPurged: StringToLineStrings(t,
					golden.FileAsString(t, "./testdata/yaml/no-importer-marker-purged.yaml")),
				Markers: map[int]*marker.Marker{},
			},
		},
		"unbalanced indentation": {
			fileName: "dummy.yaml",
			input: strings.NewReader(`
data:
  # == imptr: some_importer / begin from: ./testdata/yaml/exporter-example.yaml#[random-data] indent: align ==
  content-to-be-purged: this will be removed
      # == imptr: some_importer / end ==
`),
			wantFile: &file.File{
				FileName: "dummy.yaml",
				ContentBefore: StringToLineStrings(t, `
data:
  # == imptr: some_importer / begin from: ./testdata/yaml/exporter-example.yaml#[random-data] indent: align ==
  content-to-be-purged: this will be removed
      # == imptr: some_importer / end ==
`),
				ContentPurged: StringToLineStrings(t, `
data:
  # == imptr: some_importer / begin from: ./testdata/yaml/exporter-example.yaml#[random-data] indent: align ==
      # == imptr: some_importer / end ==
`),
				Markers: map[int]*marker.Marker{
					3: {
						Name:               "some_importer",
						LineToInsertAt:     3,
						TargetPath:         "./testdata/yaml/exporter-example.yaml",
						TargetExportMarker: "random-data",
						Indentation: &marker.Indentation{
							Mode:              marker.AlignIndentation,
							MarkerIndentation: 2, // not 6
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

func TestParseYAMLFail(t *testing.T) {
	cases := map[string]struct {
		// Input
		fileName string
		input    io.Reader

		// Output
		wantErr error
	}{
		"importer option missing": {
			fileName: "dummy.yaml",
			input: strings.NewReader(`
data:
  # == imptr: some_importer / begin ==
  content-to-be-purged: this will be removed
  # == imptr: some_importer / end ==
`),
			wantErr: marker.ErrInvalidSyntax,
		},
		"file line range not number - lower bound": {
			fileName: "dummy.yaml",
			input: strings.NewReader(`
data:
  # == imptr: some_importer / begin from: ./somefile#NOT_NUMBER~2233 ==
  content-to-be-purged: this will be removed
  # == imptr: some_importer / end == -->
`),
			wantErr: marker.ErrInvalidSyntax,
		},
		"file line range not number - upper bound": {
			fileName: "dummy.yaml",
			input: strings.NewReader(`
data:
  # == imptr: some_importer / begin from: ./somefile#1~NOT_NUMBER ==
  content-to-be-purged: this will be removed
  # == imptr: some_importer / end ==
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

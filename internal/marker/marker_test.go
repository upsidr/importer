package marker_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/upsidr/importer/internal/marker"
)

func TestNewMarker(t *testing.T) {
	cases := map[string]struct {
		input *marker.RawMarker

		want *marker.Marker
	}{
		"Line range": {
			input: &marker.RawMarker{
				Name:           "simple-marker",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 3,
				Options:        "from: ./abc.md#3~5",
			},
			want: &marker.Marker{
				Name:               "simple-marker",
				LineToInsertAt:     3,
				TargetPath:         "./abc.md",
				TargetExportMarker: "",
				TargetLines:        nil,
				TargetLineFrom:     3,
				TargetLineTo:       5,
				Indentation:        nil,
			},
		},
		"Line array": {
			input: &marker.RawMarker{
				Name:           "simple-marker",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 3,
				Options:        "from: ./abc.md#3,4,5,6",
			},
			want: &marker.Marker{
				Name:               "simple-marker",
				LineToInsertAt:     3,
				TargetPath:         "./abc.md",
				TargetExportMarker: "",
				TargetLines:        []int{3, 4, 5, 6},
				Indentation:        nil,
			},
		},
		"Exporter": {
			input: &marker.RawMarker{
				Name:           "simple-marker",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 3,
				Options:        "from: ./abc.md#[from-exporter-marker]",
			},
			want: &marker.Marker{
				Name:               "simple-marker",
				LineToInsertAt:     3,
				TargetPath:         "./abc.md",
				TargetExportMarker: "from-exporter-marker",
			},
		},
		"Exporter with absolute indent": {
			input: &marker.RawMarker{
				Name:           "simple-marker",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 3,
				Options:        "from: ./abc.md#[from-exporter-marker] indent: absolute 2",
			},
			want: &marker.Marker{
				Name:               "simple-marker",
				LineToInsertAt:     3,
				TargetPath:         "./abc.md",
				TargetExportMarker: "from-exporter-marker",
				Indentation: &marker.Indentation{
					Mode:   marker.AbsoluteIndentation,
					Length: 2,
				},
			},
		},
		"Exporter with extra indent": {
			input: &marker.RawMarker{
				Name:           "simple-marker",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 3,
				Options:        "from: ./abc.md#[from-exporter-marker] indent: extra 4",
			},
			want: &marker.Marker{
				Name:               "simple-marker",
				LineToInsertAt:     3,
				TargetPath:         "./abc.md",
				TargetExportMarker: "from-exporter-marker",
				Indentation: &marker.Indentation{
					Mode:   marker.ExtraIndentation,
					Length: 4,
				},
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := marker.NewMarker(tc.input)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("prepend result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestNewMarkerFail(t *testing.T) {
	cases := map[string]struct {
		input *marker.RawMarker

		wantErr error
	}{
		"Name missing": {
			input: &marker.RawMarker{
				Name:           "", // important
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 1,
				Options:        "dummy",
			},
			wantErr: marker.ErrMissingName,
		},
		"Marker missing matching begin and end": {
			input: &marker.RawMarker{
				Name:           "dummy",
				IsBeginFound:   true,
				IsEndFound:     false, // important
				LineToInsertAt: 1,
				Options:        "dummy",
			},
			wantErr: marker.ErrNoMatchingMarker,
		},
		"Invalid input for options": {
			input: &marker.RawMarker{
				Name:           "dummy",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 1,
				Options:        "from: ",
			},
			wantErr: marker.ErrInvalidSyntax,
		},
		"Invalid input for line numbers": {
			input: &marker.RawMarker{
				Name:           "dummy",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 1,
				Options:        "from: ./abc.md#3.5",
			},
			wantErr: marker.ErrInvalidSyntax,
		},
		"Invalid input for line range: upper bound is not a number": {
			input: &marker.RawMarker{
				Name:           "dummy",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 1,
				Options:        "from: ./abc.md#3~xyz",
			},
			wantErr: marker.ErrInvalidSyntax,
		},
		"Invalid input for line range: lower bound is not a number": {
			input: &marker.RawMarker{
				Name:           "dummy",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 1,
				Options:        "from: ./abc.md#abc~5",
			},
			wantErr: marker.ErrInvalidSyntax,
		},
		"Invalid input for filename": {
			input: &marker.RawMarker{
				Name:           "dummy",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 1,
				Options:        "from: ./some_dir/#3~5",
			},
			wantErr: marker.ErrInvalidPath,
		},
		"Invalid input for indentation": {
			input: &marker.RawMarker{
				Name:           "dummy",
				IsBeginFound:   true,
				IsEndFound:     true,
				LineToInsertAt: 1,
				Options:        "from: ./xyz.yaml#3 indent: absolute 999999999999999999999",
			},
			wantErr: marker.ErrInvalidSyntax,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := marker.NewMarker(tc.input)
			if err == nil {
				t.Fatal("error was expected but got none")
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
			}
		})
	}
}

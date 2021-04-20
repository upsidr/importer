package file

import (
	"errors"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/testingutil/golden"
)

func TestReplaceWithAfter(t *testing.T) {
	cases := map[string]struct {
		input *File

		// Output
		want    string
		wantErr error
	}{
		"test": {
			input: &File{
				FileName: "tmpfile.txt",
				ContentBefore: []string{
					"Some data",
					"and more",
				},
				ContentAfter: []byte(`Completely different data`),
			},
			want: "Completely different data",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.input.ReplaceWithAfter()
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}
			defer os.Remove(tc.input.FileName)

			result := golden.FileAsString(t, tc.input.FileName)
			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestReplaceWithPurged(t *testing.T) {
	cases := map[string]struct {
		input *File

		// Output
		want    string
		wantErr error
	}{
		"test": {
			input: &File{
				FileName: "tmpfile.txt",
				ContentBefore: []string{
					"Some data",
					"and more",
				},
				ContentPurged: []string{
					"Some purged data",
				},
				ContentAfter: []byte(`Completely different data`),
			},
			want: "Some purged data\n",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := tc.input.ReplaceWithPurged()
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}
			defer os.Remove(tc.input.FileName)

			result := golden.FileAsString(t, tc.input.FileName)
			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

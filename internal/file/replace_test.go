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
		input   *File
		options []ReplaceOption

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
		"dry-run": {
			input: &File{
				FileName: "tmpfile.txt",
				ContentBefore: []string{
					"Some data",
					"and more",
				},
				ContentAfter: []byte(`Completely different data`),
			},
			options: []ReplaceOption{
				WithDryRun(),
			},
			want: "", // not written
		},
		"skip-update found": {
			input: &File{
				FileName: "tmpfile.txt",
				ContentBefore: []string{
					"Some data",
					"and more",
				},
				ContentAfter: []byte(`Completely different data`),
				SkipUpdate:   true,
			},
			want: "", // not written
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := os.Create(tc.input.FileName)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tc.input.FileName)

			err = tc.input.ReplaceWithAfter(tc.options...)
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}

			result := golden.FileAsString(t, tc.input.FileName)
			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestReplaceWithPurged(t *testing.T) {
	cases := map[string]struct {
		input   *File
		options []ReplaceOption

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

		"dry-run": {
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
			options: []ReplaceOption{
				WithDryRun(),
			},
			want: "", // not written
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := os.Create(tc.input.FileName)
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(tc.input.FileName)

			err = tc.input.ReplaceWithPurged(tc.options...)
			if err != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
				}
				return
			}

			result := golden.FileAsString(t, tc.input.FileName)
			if diff := cmp.Diff(tc.want, result); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestReplaceFail(t *testing.T) {
	tmp := "test_file_for_replace_fail"

	f, err := os.Create(tmp)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp)

	// Not writeable
	f.Chmod(0444)

	err = replace(tmp, []byte(`some data`), &replaceMode{})
	if err == nil {
		t.Fatal("should be permission error, but got no error")
	}
}

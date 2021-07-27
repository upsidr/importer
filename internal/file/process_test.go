package file

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/marker"
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
				Markers: map[int]*marker.Marker{
					2: {
						Name:           "test annotation",
						LineToInsertAt: 2,
						TargetPath:     "../../testdata/other/note.txt",
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
				Markers: map[int]*marker.Marker{
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

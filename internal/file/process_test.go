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
						ImportTargetFile: marker.ImportTargetFile{
							Type: marker.PathBased,
							File: "../../testdata/other/note.txt",
						},
						ImportLogic: marker.ImportLogic{
							Type:  marker.CommaSeparatedLines,
							Lines: []int{1, 2},
						},
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
						ImportTargetFile: marker.ImportTargetFile{
							Type: marker.PathBased,
							File: "../../does-not-exist.txt",
						},
						ImportLogic: marker.ImportLogic{
							Type:  marker.CommaSeparatedLines,
							Lines: []int{1, 2},
						},
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

func TestRemoveMarkers(t *testing.T) {
	cases := map[string]struct {
		// Input
		file *File

		// Output
		want []byte
	}{
		"markdown: without importer marker": {
			file: &File{
				FileName: "test-file.md",
				ContentAfter: []byte(`
A
B
C
D
E

1
2
3
4
5
`),
			},
			want: []byte(`
A
B
C
D
E

1
2
3
4
5
`),
		},
		"markdown: with importer marker, should be removed": {
			file: &File{
				FileName: "test-file.md",
				ContentAfter: []byte(`
<!-- == importer: abc / begin from: abc.md#1~ == -->
A
B
C
D
E
<!-- == importer: abc / end == -->

<!-- == export: bbb / begin == -->
1
2
3
4
5
<!-- == export: bbb / end == -->
`),
			},
			want: []byte(`
A
B
C
D
E

1
2
3
4
5
`),
		},
		"yaml: with importer marker, should be removed": {
			file: &File{
				FileName: "test-file.yaml",
				ContentAfter: []byte(`
# == i: abc / begin from: some-file.yaml#[abc] ==
a:
  b:
    c: data
# == i: abc / end ==
`),
			},
			want: []byte(`
a:
  b:
    c: data
`),
		},
		"unknown file type: keep input as is": {
			file: &File{
				FileName: "test-file.dummy",
				ContentAfter: []byte(`
YAML like file, but not considered yaml due to the file extension.

# == i: abc / begin from: some-file.yaml#[abc] ==
a:
  b:
    c: data
# == i: abc / end ==
`),
			},
			want: []byte(`
YAML like file, but not considered yaml due to the file extension.

# == i: abc / begin from: some-file.yaml#[abc] ==
a:
  b:
    c: data
# == i: abc / end ==
`),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tc.file.RemoveMarkers()
			if diff := cmp.Diff(tc.want, tc.file.ContentAfter); diff != "" {
				t.Errorf("parsed result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

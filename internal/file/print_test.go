package file

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/testingutil/stdout"
)

func TestPrintAfter(t *testing.T) {
	cases := map[string]struct {
		input *File
		want  string
	}{
		"simple test": {
			input: &File{
				ContentAfter: []byte(`data
some more data
`),
			},
			want: `data
some more data
`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fakeStdout := stdout.New(t)
			defer fakeStdout.Close()

			if err := tc.input.PrintAfter(); err != nil {
				t.Fatalf("unexpected error, %v", err)
			}

			stdout := fakeStdout.ReadAllAndClose(t)

			if diff := cmp.Diff(tc.want, string(stdout)); diff != "" {
				t.Errorf("printed data didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestPrintBefore(t *testing.T) {
	cases := map[string]struct {
		input *File
		want  string
	}{
		"simple test": {
			input: &File{
				ContentBefore: []string{"data", "another data"},
			},
			want: `data
another data`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fakeStdout := stdout.New(t)
			defer fakeStdout.Close()

			if err := tc.input.PrintBefore(); err != nil {
				t.Fatalf("unexpected error, %v", err)
			}

			stdout := fakeStdout.ReadAllAndClose(t)

			if diff := cmp.Diff(tc.want, string(stdout)); diff != "" {
				t.Errorf("printed data didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestPrintPurged(t *testing.T) {
	cases := map[string]struct {
		input *File
		want  string
	}{
		"simple test": {
			input: &File{
				ContentPurged: []string{"data", "another data"},
			},
			want: `data
another data`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fakeStdout := stdout.New(t)
			defer fakeStdout.Close()

			if err := tc.input.PrintPurged(); err != nil {
				t.Fatalf("unexpected error, %v", err)
			}

			stdout := fakeStdout.ReadAllAndClose(t)

			if diff := cmp.Diff(tc.want, string(stdout)); diff != "" {
				t.Errorf("printed data didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

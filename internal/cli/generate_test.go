package cli

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/parse"
	"github.com/upsidr/importer/internal/testingutil/golden"
	"github.com/upsidr/importer/internal/testingutil/stdout"
)

// Run `go test ./... -updateGolden` to updateGolden golden files under testdata
var updateGolden = flag.Bool("update", false, "update golden files")

func TestGenerateStdout(t *testing.T) {
	cases := map[string]struct {
		// Input
		inputFile string

		// Output
		wantFile      string
		wantErrString string
	}{
		"markdown": {
			inputFile: "../../testdata/markdown/simple-before.md",
			wantFile:  "../../testdata/markdown/simple-updated.md",
		},
		"markdown with exporter": {
			inputFile: "../../testdata/markdown/using-exporter-before.md",
			wantFile:  "../../testdata/markdown/using-exporter-updated.md",
		},
		"error case: file not found": {
			inputFile:     "does_not_exist",
			wantErrString: "no such file",
		},
		"error case: file not supported (.txt)": {
			inputFile:     "../../testdata/other/note.txt",
			wantErrString: parse.ErrUnsupportedFileType.Error(),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			fakeStdout := stdout.New(t)
			defer fakeStdout.Close()

			err := generate(tc.inputFile, "") // Empty second argument means generate writes to stdout
			if err != nil {
				if !strings.Contains(err.Error(), tc.wantErrString) {
					t.Fatalf("error with generate, %v", err)
				}
				return
			}

			stdout := fakeStdout.ReadAllAndClose(t)

			got := string(stdout)
			want := golden.FileAsString(t, tc.wantFile)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestGenerateToFile(t *testing.T) {
	cases := map[string]struct {
		// Input
		inputFile string

		// Output
		wantFile string
	}{
		"markdown": {
			inputFile: "../../testdata/markdown/simple-before.md",
			wantFile:  "../../testdata/markdown/simple-updated.md",
		},
		"markdown with exporter": {
			inputFile: "../../testdata/markdown/using-exporter-before.md",
			wantFile:  "../../testdata/markdown/using-exporter-updated.md",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tempFile, err := os.CreateTemp("", "")
			if err != nil {
				t.Fatal(err)
			}

			err = generate(tc.inputFile, tempFile.Name()) // Second argument for target file
			if err != nil {
				t.Fatalf("error with generate, %v", err)
			}

			if *updateGolden {
				processed := golden.File(t, tempFile.Name())
				golden.UpdateFile(t, tc.wantFile, processed)
			}

			got := golden.FileAsString(t, tempFile.Name())
			want := golden.FileAsString(t, tc.wantFile)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

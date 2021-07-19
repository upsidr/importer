package cli

import (
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/testingutil/golden"
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
			inputFile: "../../testdata/simple-before.md",
			wantFile:  "../../testdata/simple-after.md",
		},
		"markdown long input": {
			inputFile: "../../testdata/long-input-purged.md",
			wantFile:  "../../testdata/long-input-after.md",
		},
		"markdown with exporter": {
			inputFile: "../../testdata/using-exporter-before.md",
			wantFile:  "../../testdata/using-exporter-after.md",
		},
		"error case: file not found": {
			inputFile:     "does_not_exist",
			wantErrString: "no such file",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			// Stdout setup
			origStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatal(err)
			}
			os.Stdout = w

			err = generate(tc.inputFile, "") // Empty second argument means generate writes to stdout
			if err != nil {
				if !strings.Contains(err.Error(), tc.wantErrString) {
					t.Fatalf("error with generate, %v", err)
				}
				return
			}

			// Get stdout back
			w.Close()
			os.Stdout = origStdout
			stdout, err := io.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}
			r.Close()

			if *updateGolden {
				golden.UpdateFile(t, tc.wantFile, stdout)
			}

			got := string(stdout)
			want := golden.FileAsString(t, tc.wantFile)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

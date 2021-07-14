package cli

import (
	"flag"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/testingutil/golden"
)

// Run `go test ./... -updateGolden` to updateGolden golden files under testdata
var updateGolden = flag.Bool("update", false, "update golden files")

func TestGenerate(t *testing.T) {
	cases := map[string]struct {
		// Input
		inputFile string

		// Output
		wantFile string
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
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			copiedFile, remove := golden.CopyTemp(t, tc.inputFile)
			defer remove()

			err := generate(copiedFile)
			if err != nil {
				t.Fatalf("error with generate, %v", err)
			}

			if *updateGolden {
				processed := golden.File(t, copiedFile)
				golden.UpdateFile(t, tc.wantFile, processed)
			}

			got := golden.FileAsString(t, copiedFile)
			want := golden.FileAsString(t, tc.wantFile)

			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

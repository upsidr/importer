package cli

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/testingutil/golden"
)

func TestUpdate(t *testing.T) {
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
			inputFile: "../../testdata/markdown/import-with-exporter-before.md",
			wantFile:  "../../testdata/markdown/import-with-exporter-updated.md",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			copiedFile, remove := golden.CopyTemp(t, tc.inputFile)
			defer remove()

			err := update(copiedFile)
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

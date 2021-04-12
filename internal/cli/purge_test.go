package cli

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/testingutil/golden"
)

func TestPurge(t *testing.T) {
	cases := map[string]struct {
		// Input
		inputFile string

		// Output
		wantFile string
	}{
		"markdown": {
			inputFile: "../../testdata/simple-before.md",
			wantFile:  "../../testdata/simple-purged.md",
		},
		"markdown long input": {
			inputFile: "../../testdata/long-input-after.md",
			wantFile:  "../../testdata/long-input-purged.md",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			copiedFile, remove := golden.CopyTemp(t, tc.inputFile)
			defer remove()

			err := purge(copiedFile)
			if err != nil {
				t.Fatalf("error with generate, %v", err)
			}

			if *update {
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

package marker

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestHandleAbsoluteIndentation(t *testing.T) {
	cases := map[string]struct {
		originalSlice        []byte
		exporterMarkerIndent int
		targetIndent         int

		want []byte
	}{
		"case 1. - original data has more indent": {
			originalSlice:        []byte("        abcdef"), // 8 spaces
			exporterMarkerIndent: 6,                        // remove 6
			targetIndent:         4,
			want:                 []byte("      abcdef"), // This is 8 - 6 + 4 = 6
		},
		"case 2. - original data has less indent": {
			originalSlice:        []byte("    abcdef"), // 4 spaces
			exporterMarkerIndent: 2,                    // remove 2
			targetIndent:         10,
			want:                 []byte("            abcdef"), // This is 4 - 2 + 10 = 12
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := handleAbsoluteIndentation(tc.originalSlice, tc.exporterMarkerIndent, tc.targetIndent)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("prepend result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestPrependWhitespaces(t *testing.T) {
	cases := map[string]struct {
		originalSlice   []byte
		whitespaceCount int

		want []byte
	}{
		"indentation": {
			originalSlice:   []byte("abcdef"),
			whitespaceCount: 6,
			want:            []byte("      abcdef"),
		},
		"extra indentation": {
			originalSlice:   []byte("  abcdef"),
			whitespaceCount: 6,
			want:            []byte("        abcdef"),
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got := prependWhitespaces(tc.originalSlice, tc.whitespaceCount)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("prepend result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

package parse

import (
	"errors"
	"io"
	"strings"
	"testing"
)

func TestParseFail(t *testing.T) {
	cases := map[string]struct {
		// Input
		fileName string
		input    io.Reader

		// Output
		wantErr error
	}{
		"extension not supported": {
			fileName: "no_extension",
			input:    strings.NewReader("dummy"),
			wantErr:  ErrUnsupportedFileType,
		},
		"no data": {
			fileName: "dummy.md",
			input:    nil,
			wantErr:  ErrNoInput,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := Parse(tc.fileName, tc.input)
			if err == nil {
				t.Fatal("error was expected but got none")
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
			}
		})
	}
}

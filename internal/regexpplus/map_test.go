package regexpplus_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/regexpplus"
)

func TestFindNamedSubgroups(t *testing.T) {
	cases := map[string]struct {
		targetLine  string
		regexpInput string

		want map[string]string
	}{
		"2 different named groups": {
			targetLine:  "abc def ghi jkl mno",
			regexpInput: `(?P<a>a.*).*(?P<d>d)`,
			want: map[string]string{
				"":  "abc d",
				"a": "abc ",
				"d": "d",
			},
		},
		"2 named groups, with duplicating name": {
			targetLine:  "abc def ghi jkl mno",
			regexpInput: `(?P<a>a).*(?P<a>d)`,
			want: map[string]string{
				"":  "abc d",
				"a": "d", // later one takes precedence
			},
		},
		"2 unnamed groups": {
			targetLine:  "abc def ghi jkl mno",
			regexpInput: `(ab).*(de)`,
			want: map[string]string{
				"": "de", // later group takes precedence
			},
		},
		"no group": {
			targetLine:  "abc def ghi jkl mno",
			regexpInput: `d.*`,
			want: map[string]string{
				"": "def ghi jkl mno",
			},
		},
		"simple input": {
			targetLine:  "abc def ghi jkl mno",
			regexpInput: `def`,
			want: map[string]string{
				"": "def",
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			got, err := regexpplus.MapWithNamedSubgroups(tc.targetLine, tc.regexpInput)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestFindNamedSubgroupsFail(t *testing.T) {
	cases := map[string]struct {
		targetLine  string
		regexpInput string

		wantErr error
	}{
		"2 different named groups": {
			targetLine:  "abc def ghi jkl mno",
			regexpInput: `xyz`,
			wantErr:     regexpplus.ErrNoMatch,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := regexpplus.MapWithNamedSubgroups(tc.targetLine, tc.regexpInput)
			if err == nil {
				t.Fatal("error was expected but got none")
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("error did not match:\n    want: %v\n    got:  %v", tc.wantErr, err)
			}
		})
	}
}

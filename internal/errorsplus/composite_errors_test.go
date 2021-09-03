package errorsplus_test

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/upsidr/importer/internal/errorsplus"
)

var (
	errFirst  = errors.New("first error")
	errSecond = errors.New("second error")
	errThird  = errors.New("third error")
	errOther  = errors.New("some other error")
)

func TestCompositeErrors(t *testing.T) {
	cases := map[string]struct {
		errs          errorsplus.Errors
		wantErr       error
		wantErrString string
	}{
		"zero error": {
			errs: errorsplus.Errors{
				// empty
			},
			wantErr:       nil, // cannot be checked against nil, this is skipped explicitly below
			wantErrString: "",  // empty
		},
		"single error": {
			errs: errorsplus.Errors{
				errFirst,
			},
			wantErr:       errFirst,
			wantErrString: "first error",
		},
		"2 errors": {
			errs: errorsplus.Errors{
				errFirst,
				errSecond,
			},
			wantErr: errFirst,
			wantErrString: `composite error:
	first error
	second error`,
		},
		"3 errors": {
			errs: errorsplus.Errors{
				errFirst,
				errSecond,
				errThird,
			},
			wantErr: errFirst,
			wantErrString: `composite error:
	first error
	second error
	third error`,
		},
		"duplicated errors": {
			errs: errorsplus.Errors{
				errFirst,
				errFirst,
				errFirst,
			},
			wantErr: errFirst,
			wantErrString: `composite error:
	first error
	first error
	first error`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if tc.wantErr != nil && !errors.Is(tc.errs, tc.wantErr) {
				t.Errorf("error mismatch\n\tcomposite error: %v\n\twanted error: %v", tc.errs, tc.wantErr)
			}

			errString := tc.errs.Error()
			if diff := cmp.Diff(tc.wantErrString, errString); diff != "" {
				t.Errorf("result didn't match (-want / +got)\n%s", diff)
			}
		})
	}
}

func TestCompositeErrorsMismatch(t *testing.T) {
	cases := map[string]struct {
		errs          errorsplus.Errors
		mismatchedErr error
	}{
		"zero error": {
			errs: errorsplus.Errors{
				// empty
			},
			mismatchedErr: errOther,
		},
		"single error": {
			errs: errorsplus.Errors{
				errFirst,
			},
			mismatchedErr: errOther,
		},
		"2 errors": {
			errs: errorsplus.Errors{
				errFirst,
				errSecond,
			},
			mismatchedErr: errOther,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if errors.Is(tc.errs, tc.mismatchedErr) {
				t.Errorf("unexpected error match, where it should not be matched with errors.Is\n\tcomposite error: %v\n\ttarget error: %v", tc.errs, tc.mismatchedErr)
			}
		})
	}
}

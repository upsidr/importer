package marker

import "errors"

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrNoInput             = errors.New("no file content found")
	ErrInvalidPath         = errors.New("invalid path provided")
	ErrInvalidSyntax       = errors.New("invalid syntax given")
	ErrNoMatchingMarker    = errors.New("no matching marker found")
	ErrMissingOption       = errors.New("missing option")
	ErrMissingName         = errors.New("unnamed marker cannot be used")
	ErrNoFileInput         = errors.New("no target file input provided")
	ErrInvalidURL          = errors.New("invalid URL")
	ErrGetMarkerTarget     = errors.New("failed to get marker target")
	ErrNonSuccessCode      = errors.New("received non-success error code")
)

package marker

import "errors"

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
	ErrNoInput             = errors.New("no file content found")
	ErrInvalidPath         = errors.New("invalid path provided")
	ErrInvalidSyntax       = errors.New("invalid syntax given")
	ErrNoMatchingMarker    = errors.New("no matching marker found, marker must be a begin/end pair")
	ErrMissingOption       = errors.New("missing option")
	ErrMissingName         = errors.New("unnamed marker cannot be used")
)

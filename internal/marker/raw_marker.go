package marker

import "fmt"

// RawMarker holds marker data as they are provided. This is only a data holder
// before processing the marker inputs to generate valid marker data. RawMarker
// can be missing some fields, which makes it possibly invalid as a Marker.
type RawMarker struct {
	Name                 string
	IsBeginFound         bool
	IsEndFound           bool
	LineToInsertAt       int
	Options              string
	PrecedingIndentation string
}

// Validate checks RawMarker's validity.
func (r *RawMarker) Validate() error {
	if r.Name == "" {
		return fmt.Errorf("%w", ErrMissingName)
	}

	if !r.IsBeginFound || !r.IsEndFound {
		return fmt.Errorf("%w", ErrNoMatchingMarker)
	}

	return nil
}

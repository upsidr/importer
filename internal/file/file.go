package file

import "github.com/upsidr/importer/internal/marker"

// File holds onto file data. This does not provide any file processing support
// by itself.
type File struct {
	FileName string

	// ContentBefore holds the file content as it was before processing. The
	// first slice represents the line number, and the second is for the actual
	// data.
	ContentBefore []string

	// ContentPurged holds the file coontent, but removes the parts between
	// importer annotation begin/end. The first slice represents the line
	// number, and the second is for the actual data.
	ContentPurged []string

	// ContentAfter holds the file content after the import has been run. This
	// only holds the actual data in byte slice representation.
	ContentAfter []byte

	// Markers is an array holding onto each annotation block.
	Markers map[int]*marker.Marker

	// SkipUpdate is used to skip updating the file in place. This is set by a
	// special marker syntax.
	SkipUpdate bool
}

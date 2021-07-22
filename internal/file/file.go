package file

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

	// Annotations is an array holding onto each annotation block.
	Annotations map[int]*Annotation
}

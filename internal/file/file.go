package file

// File holds onto file data.
type File struct {
	FileName string

	// fileType is derived from FileName, which is simply represented using
	// extension format.
	fileType string

	// contentBefore holds the file content as it was before processing. The
	// first slice represents the line number, and the second is for the actual
	// data.
	contentBefore [][]byte

	// contentPurged holds the file coontent, but removes the parts between
	// importer annotation begin/end. The first slice represents the line
	// number, and the second is for the actual data.
	contentPurged [][]byte

	// contentAfter holds the file content after the import has been run. This
	// only holds the actual data in byte slice representation.
	contentAfter []byte

	// annotations is an array holding onto each annotation block.
	annotations map[int]annotation
}

type annotation struct {
	name                  string
	lineWithBeginOriginal int
	lineWithBeginPurged   int
	targetPath            string
	targetLines           []int
}

func NewFile() *File {
	result := &File{}

	return result
}

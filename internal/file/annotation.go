package file

// Annotation holds on to the data required for importer processing. This does
// not hold the target file content itself, and that needs to be handled
// separately.
type Annotation struct {
	Name           string
	LineToInsertAt int

	TargetPath string

	TargetExportMarker string

	TargetLines    []int
	TargetLineFrom int
	TargetLineTo   int

	// TODO: Add insert style such as code verbatim, details, quotes, etc.
}

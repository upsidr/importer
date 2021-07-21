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

	ExtraIndentation    int
	AbsoluteIndentation int
	// TODO: Add insert style such as code verbatim, details, quotes, etc.
}

var (
	// ExportMarkerMarkdown is the marker used to indicate how a file can
	// export specific sections.
	//
	// Example:
	//   <!-- == export: simple_instruction / begin == -->
	//   This is the content that will be exported under "simple_instruction" name.
	//   You can import this content by providing option such as:
	//     ./file_path.txt#[simple_instruction]
	//   <!-- == export: simple_instruction / end == -->
	ExportMarkerMarkdown = `<!-- == export: (?P<export_marker_name>\S+) \/ (?P<exporter_marker_condition>begin|end) == -->`

	// ExportMarkerYAML is the marker used to indicate how a file can export
	// specific sections.
	//
	// Example:
	//   data:
	//     some-data: something
	//     # == export: random_data / begin ==
	//     random-data: this is exported
	//     # == export: random_data / end ==
	ExportMarkerYAML = `(?P<export_marker_indent>\s*)# == export: (?P<export_marker_name>\S+) \/ (?P<exporter_marker_condition>begin|end) ==`
)

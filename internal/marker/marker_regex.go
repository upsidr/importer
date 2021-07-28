package marker

var (
	// ImporterMarkerMarkdown is the annotation used for importer to find match.
	//
	// Example:
	//   <!-- == imptr: some_importer_name / begin from: ./file.txt#2~22 == -->
	//   <!-- == imptr: some_importer_name / end == -->
	ImporterMarkerMarkdown = `<!-- == imptr: (?P<importer_name>\S+) \/ (?P<importer_marker>begin|end)(?P<importer_option>.*) == -->`

	// ImporterMarkerYAML is the annotation used for importer to find match.
	ImporterMarkerYAML = `# == imptr: (?P<importer_name>\S+) \/ (?P<importer_marker>begin|end)(?P<importer_option>.*) ==`

	// OptionFilePathIndicator is the pattern used for parsing Importer file options.
	OptionFilePathIndicator = `from: (?P<importer_target_path>\S+)\s*\#(?P<importer_target_detail>[0-9a-zA-Z,-_\~]+)\s?`

	// OptionIndentMode is the pattern used for specifying indentation mode.
	OptionIndentMode = `indent: (?P<importer_indent_mode>absolute|extra) (?P<importer_indent_length>\d+)`
)

var (
	// ExporterMarkerMarkdown is the marker used to indicate how a file can
	// export specific sections.
	//
	// Example:
	//   <!-- == export: simple_instruction / begin == -->
	//   This is the content that will be exported under "simple_instruction" name.
	//   You can import this content by providing option such as:
	//     ./file_path.txt#[simple_instruction]
	//   <!-- == export: simple_instruction / end == -->
	ExporterMarkerMarkdown = `<!-- == export: (?P<export_marker_name>\S+) \/ (?P<exporter_marker_condition>begin|end) == -->`

	// ExporterMarkerYAML is the marker used to indicate how a file can export
	// specific sections.
	//
	// Example:
	//   data:
	//     some-data: something
	//     # == export: random_data / begin ==
	//     random-data: this is exported
	//     # == export: random_data / end ==
	ExporterMarkerYAML = `(?P<export_marker_indent>\s*)# == export: (?P<export_marker_name>\S+) \/ (?P<exporter_marker_condition>begin|end) ==`
)
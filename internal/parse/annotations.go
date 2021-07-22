package parse

var (
	// ImporterAnnotationMarkdown is the annotation used for importer to find match.
	//
	// Example:
	//   <!-- == imptr: some_importer_name / begin from: ./file.txt#2~22 == -->
	//   <!-- == imptr: some_importer_name / end == -->
	ImporterAnnotationMarkdown = `<!-- == imptr: (?P<importer_name>\S+) \/ (?P<importer_marker>begin|end)(?P<importer_option>.*) == -->`

	// ImporterAnnotationYAML is the annotation used for importer to find match.
	ImporterAnnotationYAML = `# == imptr: (?P<importer_name>\S+) \/ (?P<importer_marker>begin|end)(?P<importer_option>.*) ==`

	// OptionFilePathIndicator is the pattern used for parsing Importer file options.
	OptionFilePathIndicator = `from: (?P<importer_target_path>\S+)\s*\#(?P<importer_target_detail>[0-9a-zA-Z,-_\~]+)\s?`

	// OptionIndentMode is the pattern used for specifying indentation mode.
	OptionIndentMode = `indent: (?P<importer_indent_mode>absolute|extra) (?P<importer_indent_length>\d+)`
)

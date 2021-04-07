package parse

var (
	// AnnotationMarkdown is the annotation used for importer to find match.
	AnnotationMarkdown = `<!-- == imptr: (?P<importer_name>\S+) \/ (?P<importer_marker>begin|end)(?P<importer_option>.*) == -->`

	// AnnotationYAML is the annotation used for importer to find match.
	AnnotationYAML = `# == imptr: (?P<importer_name>\S+) \/ (?P<importer_marker>begin|end)(?P<importer_option>.*) ==`

	// FilePathIndicator
	OptionFilePathIndicator = `from: (?P<importer_target_path>\S+)\s*\#(?P<importer_target_lines>\S+)\s?`
)

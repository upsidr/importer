## Roadmap

The items here are expected to be complete before v1.0 release. The items are not in priority order, though the top most ones tend to be tackled earlier.

### Support YAML files

YAML is another file type that needs to be a single file input, and thus you cannot pull in other content.

Importer is planning to support YAML files in the near future.

### Add `--dry-run` flag

`importer preview` does very basic preview of how the file would be updated. This should be updated so that when running Importer command with flag `--dry-run` would get the output to stdout.

### Add `generate` and `update` commands

Currently, `importer generate` takes in a file argument, and updates the file content in place.

Instead, we are aiming to provide `importer update` to provide the same feature, while `importer generate` to output the result to stdout. This allows having a separate file that contains Importer Annotations, and a generated file as a separate file.

### Add `graph` command

Currently, Importer only looks at the provided argument and its Import Target Files. When the Target File contains another Importer Annotation, it would be better to update the Target File content first.
We will need much better processing than simple regex handling, and abstract syntax tree needs to be created for this command. Also, this command needs to ensure there is no cyclic dependencies in the Importer definitions.

### Support line brak in Importer Annotation and Export Marker

Currently Importer Annotation and Export Marker have to be a single line input. If you have a line break in them, it will be ignored. This is because how it's currently implemented, and fixing this would require a proper AST setup when parsing a file.

### Add special annotations `ignore` to skip Importer run

When having an automation such as `find . -name '*.md' -exec importer generate {} \;`, you may want to skip some files.

This shouldn't skip Export Marker handling, though.

### Add `diff` command

Provide a nice diff command where you can see how Importer changes the file content.

### Support pulling files from internet

Just like `kubectl`, support providing a URL for the Import Target.

## Potential Items

The items here are being considered at the moment, but there is no clear timeline. They need more input as they seem to help for some cases, but may not be too useful for many.

### Add Importer Config - To be confirmed

Importer handles the target files by relative paths, but we may want to support absolute path. In order to do that, though, we may need to have some separate configuration at the root of repository (in case of using Git repo), and use that location as the root. There could be some other benefits for having a dedicated config, but needs further consideration.

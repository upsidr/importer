# Details of Importer

## Background

Importer aims to achieve one simple goal: break long files into small files.

In most of programming languages, you can define variables and reuse code, and some even try to emphasise on DRY - "Don't Repeat Yourself". Importer's goal is similar in idea, but it tries to be as "dumb" as possible. It doesn't guarantee the code reuse is done all the time, there is essentially no compilation, and no complicated logic at all. This means, if the file gets manually updated after Importer update, Importer doesn't care. It is a simple enough tool, which can be embedded as a part of some other automation. If you want to ensure the Importer's `update` command is always run against a given file, you should be putting some CI job using Importer - but Importer shouldn't know anything about CI itself.

Also, Importer's idea is not about DRY - instead, it's about "reuse for clarity".

For example, in any programming language, if you update a single variable, it could have a huge impact as the variable could be referred in many places. You will need some smart tools such as IDE, CI jobs, or test caess to fully understand the impact. You may miss some change that wasn't intended to take place because of the hidden dependencies.

However, you probably won't simply copy and paste the same code in multiple places. It may be very easy to see even for reviewer, but when it comes to updating them, you would have to update in many places. This is very error prone, as you can forget to update one file, and it is often difficult to spot those.

That's why Importer was created. You can stick to simple file format, but when it becomes harder to maintain due to the size and complexity, you can split up the file into multiple files and "Import" as necessary. This keeps the actual file as is, meaning Markdown will still be a pure Markdown.

## Implementation Details

Importer is a simple regular expression file reader. It looks for special Importer Marker comment, which has no meaning in that given language, but Importer can parse such comment and wires up other file content. [You can find more about Importer Marker here.](../getting-started/annotatinos.md)

In language like Markdown and YAML, a single file is made to be as is, meaning you cannot import other files. This makes them really simple and easy to get started, but when you try to do something a bit more involved, it becomes difficult to maintain very quickly.

Because Importer tries to be "dumb", it doesn't actually know much about the given file syntax. Importer looks for Importer Marker comments, parses them, and generates the updated version of that file, with specified lines imported into it.

Because the goal of Importer is very simple, the implementation is based on simple regular expressions. It is not made to be performant, nor capable of handling complex scenarios. But it works for most cases, such as Markdown and YAML. Other file typse may benefit from this approach. If there is any other file types that could benefit from this, we will look to expand our support in the future.

<!-- == imptr: roadmap / begin from: ./roadmap.md#1~ == -->

# Roadmap

The items here are expected to be complete before v1.0 release. The items are not in priority order, though the top most ones tend to be tackled earlier.

### Support YAML files

YAML is another file type that needs to be a single file input, and thus you cannot pull in other content.

Importer is planning to support YAML files in the near future.

### Add `--dry-run` flag

`importer preview` does very basic preview of how the file would be updated. This should be updated so that when running Importer command with flag `--dry-run` would get the output to stdout.

### Add `generate` and `update` commands

Currently, `importer generate` takes in a file argument, and updates the file content in place.

Instead, we are aiming to provide `importer update` to provide the same feature, while `importer generate` to output the result to stdout. This allows having a separate file that contains Importer Markers, and a generated file as a separate file.

### Add `graph` command

Currently, Importer only looks at the provided argument and its Import Target Files. When the Target File contains another Importer Marker, it would be better to update the Target File content first.
We will need much better processing than simple regex handling, and abstract syntax tree needs to be created for this command. Also, this command needs to ensure there is no cyclic dependencies in the Importer definitions.

### Support line brak in Importer Marker and Exporter Marker

Currently Importer Marker and Exporter Marker have to be a single line input. If you have a line break in them, it will be ignored. This is because how it's currently implemented, and fixing this would require a proper AST setup when parsing a file.

### Add special markers `ignore` to skip Importer run

When having an automation such as `find . -name '*.md' -exec importer generate {} \;`, you may want to skip some files.

This shouldn't skip Exporter Marker handling, though.

### Add `diff` command

Provide a nice diff command where you can see how Importer changes the file content.

### Support pulling files from internet

Just like `kubectl`, support providing a URL for the Import Target.

## Potential Items

The items here are being considered at the moment, but there is no clear timeline. They need more input as they seem to help for some cases, but may not be too useful for many.

### Add Importer Config - To be confirmed

Importer handles the target files by relative paths, but we may want to support absolute path. In order to do that, though, we may need to have some separate configuration at the root of repository (in case of using Git repo), and use that location as the root. There could be some other benefits for having a dedicated config, but needs further consideration.

<!-- == imptr: roadmap / end == -->

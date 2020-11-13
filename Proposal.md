# Design Document

Written by: @rytswd

## What it does

- In any file, you can import other file content
- Works like code generation, but does not need to be code

### Use cases

There are 2 cases this can definitely help.

- Markdown templating  
  Although markdown is only meant to be a simple document, it's used in ubiquitous areas in GitHub. It is certainly easy to get started with, but it can quickly become lengthy and content heavy. Instead of having everything written in one file, it would be great if we can import other file content into markdown. That way, the documentation can be broken down into small chunks, and can be reused in multiple places.

- YAML templating  
  Helm and Kustomize shine in this area, but if you need the most simplistic approach of using the same YAML stanza imported, Helm and Kustomize do solve the problem, but are a bit of an overkill.

Both of them are interesting because Markdown and YAML cannot import other file, because of its simplicity.

## What it does NOT

- It does not guarantee the content validity
- It is not meant to be used for many file types

## Name

Temporarily set to `importer`, which explains what it does well enough, in my opinion.

## Implementation Ideas

- Compile into a binary `importer` as CLI tool
- It takes single file as argument
- It checks input file, and looks for annotation such as `<!--IMPORTER:BEGIN - ./somefile.go#110~120 -->`
- Append the content of import reference
- Use cobra
- Create subcommands of `generate`, `preview`, and `graph`

## Release Plans

### v0.1

- Support `generate` subcommand, which creates the new file based on the referenced file embedded in the file
- Support `preview` subcommand, which outputs the embedded content together to stdout
- Support only markdown
- No nesting support (i.e. if `b.md` imports `c.md`, and you try to use `a.md` to import `b.md`, it does not reconcile `c.md` - simply takes the file as is)

### v0.2

- Create GitHub Actions integration
- Create documentation for GitHub Actions setup

### v0.3

- Support YAML

### v0.4

- Support nesting
- Support `graph` subcommand to generate dependency graph

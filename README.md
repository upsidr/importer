# Importer

Import any lines, from anywhere

[![Build Status](https://github.com/upsidr/importer/workflows/Build%20Importer/badge.svg?event=push)](build-status) | [![GitHub Release Date](https://img.shields.io/github/release-date/upsidr/importer?color=powderblue)](releases)

[build-status]: https://github.com/upsidr/importer/actions
[releases]: https://github.com/upsidr/importer/releases

![Demo](/assets/images/importer-update-demo.gif)

## 🌄 What is Importer?

Importer is a CLI tool to allow any file to import other file content, including Markdown, YAML, to name a few. Importer uses **Importer Markers**, which are often provided as comment, to find the relevant file and import defined lines based on line numbers and other details.

Files such as Markdown and YAML which are meant to be a single file input can have other file to be pulled in. Importer aims to provide this extra feature without breaking any language, and that means Importer uses code generation approach, where the **Markers** are used to update the file in place.

This may seem like an unnecessary layer for simple files such as Markdown and YAML, but this allows better structure and code reuse, while retaining or even enhancing code readability.

You can find more about the details of Importer design [here](/docs/details/details.md).

## ✨ Install

<!-- == imptr: getting-started-install / begin from: ./docs/getting-started/install.md#[homebrew-install] == -->

You can get Importer with simple Homebrew command.

```bash
$ brew install upsidr/tap/importer
```

You can also find the relevent binary files under [releases](https://github.com/upsidr/importer/releases).

<!-- == imptr: getting-started-install / end == -->

<details>
<summary>Other Installation Options</summary>

### Install with Go

<!-- == imptr: install-with-go / begin from: ./docs/getting-started/install.md#[go-get] == -->

You can also use Go to install.

```bash
$ go get github.com/upsidr/importer/cmd/importer@v0.0.1-rc2
```

<!-- == imptr: install-with-go / end == -->

</details>

## 🎮 Commands

<!-- == imptr: commands / begin from: ./docs/details/commands.md#[help-output] == -->

```console
$ importer -h
Import any lines, from anywhere

Usage:
  importer [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  generate    Processes Importer markers and send output to stdout or file
  help        Help about any command
  preview     Shows a preview of Importer update and purge results
  purge       Removes all imported lines and update the file in place
  update      Processes Importer markers and update the file in place

Flags:
  -h, --help   help for importer

Use "importer [command] --help" for more information about a command.
```

<!-- == imptr: commands / end == -->

## 🧩 Supported Files

<!-- == imptr: supported-files / begin from: ./docs/details/supported-files.md#[list] == -->

| File Type | Is Supported? | File Extensions | Additional Importer Option |
| --------- | :-----------: | --------------- | -------------------------- |
| Markdown  |      ✅       | `.md`           |                            |
| YAML      |      ✅       | `.yaml`, `.yml` | Indentation                |
| HTML      |      🚧       | TBC             |                            |
| TOML      |      🚧       | TBC             |                            |

Any other file type not specified above are not supported by Importer at the moment.

To request additional file support, please file an issue from [here](https://github.com/upsidr/importer/issues/new?assignees=&labels=enhancement&template=feature-request.yaml&title=%5BFeature+Request%5D%3A+).

<!-- == imptr: supported-files / end == -->

## 🖋 Markers

<!-- == imptr: basic-marker / begin from: ./docs/details/markers.md#[basic-marker] == -->

**Markers** are a simple comment with special syntax Importer understands. There
are several types of markers.

![Marker explained][marker-explanation]

[marker-explanation]: /assets/images/marker-explanation.png "Marker Explanation"

Importer has the following markers.

| Name                 | Description                                                     |
| -------------------- | --------------------------------------------------------------- |
| Importer Marker      | The main marker, used to import data from other file.           |
| Exporter Marker      | Supplemental markers used to define line range in target files. |
| Skip Importer Update | Special marker to suppress `importer update`.                   |
| Auto Generated Note  | Special marker for `importer generate` information.             |

<!-- == imptr: basic-marker / end == -->

You can find more about the Importer Marker [here](/docs/details/markers.md).

## 🚀 Examples

<!-- == imptr: getting-started-example-short / begin from: ./docs/getting-started/examples-markdown.md#[preview] == -->

`importer preview` command gives you a quick look at how the file may change when `importer update` and `importer purge` are run against the provided file. This is meant to be useful for testing and debugging.

```console
$ importer preview ./testdata/markdown/demo-before.md
---------------------------------------
Content Before:
1:      # Markdown Demo
2:
3:      <!-- == imptr: short-description / begin from: ./snippet-description.md#[for-demo] == -->
4:      Any content here will be replaced by Importer.
5:      <!-- == imptr: short-description / end == -->
---------------------------------------

---------------------------------------
Content After Purged:
1:      # Markdown Demo
2:
3:      <!-- == imptr: short-description / begin from: ./snippet-description.md#[for-demo] == -->
4:      <!-- == imptr: short-description / end == -->
---------------------------------------

---------------------------------------
Content After Processed:
1:      # Markdown Demo
2:
3:      <!-- == imptr: short-description / begin from: ./snippet-description.md#[for-demo] == -->
4:      This demonstrates how a markdown can import other file content.
5:
6:      Importer is a CLI tool to read and process Importer and Exporter markers.
7:      This can be easily integrated into CI/CD and automation setup.
8:      <!-- == imptr: short-description / end == -->
---------------------------------------

You can replace the file content with either of the commands below:

  importer update ./testdata/markdown/demo-before.md     Replace the file content with the Importer processed file.
  importer purge ./testdata/markdown/demo-before.md      Replace the file content by removing all data between marker pairs.

You can find more with 'importer help'
```

<!-- == imptr: getting-started-example-short / end == -->

You can find more examples [here](/docs/getting-started/examples-markdown.md).

## :octocat: GitHub Action Integration

<!-- == imptr: getting-started-github-action / begin from: ./docs/getting-started/github-actions.md#[with-homebrew] == -->

Because you can install Importer using Homebrew, you can set up GitHub Action definition such as below:

<!--TODO: The below YAML is exactly where Importer should be able to pull in the actual file content-->

```yaml
jobs:
  importer:
    name: Run Importer Generate
    runs-on: ubuntu-latest
    steps:
      - name: Install Importer
        run: brew install upsidr/tap/importer

      - name: Check out
        uses: actions/checkout@v2
        with:
          fetch-depth: 0
      - name: Run Importer against README.md
        run: importer generate README.md
      - name: Check if README.md has any change compared to the branch
        run: |
          git status --short
          git diff-index --quiet HEAD
```

This repository uses Importer to generate some of the markdown documentation.

You can find actually running CI setup in [`.github/workflows/importer-ci.yaml`](https://github.com/upsidr/importer/blob/main/.github/workflows/importer-ci.yaml).

<!-- == imptr: getting-started-github-action / end == -->

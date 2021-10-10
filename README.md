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

![Marker in Action][marker-in-action]

[marker-in-action]: /assets/images/importer-overview.png "Marker in Action"

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
$ go get github.com/upsidr/importer/cmd/importer@v0.1.0
```

<!-- == imptr: install-with-go / end == -->

</details>

## 🎮 Commands

<!-- == imptr: commands / begin from: ./docs/details/commands.md#[list] == -->

| Name                     | Description                                                                                       |
| ------------------------ | ------------------------------------------------------------------------------------------------- |
| `importer generate FILE` | Run Importer processing on `FILE`, and write the result to stdout.                                |
| `importer update FILE`   | Run Importer processing on `FILE`, and update it in place.                                        |
| `importer purge FILE`    | Parse Importer Markers, remove any content within Importer Markers, and update the file in plcae. |
| `importer preview FILE`  | Write before/purged/after preview of how Importer processes the file content to stdout.           |

<!-- == imptr: commands / end == -->

You can find more about the commands [here](/docs/details/commands.md).

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

**Markers** are a simple comment with special syntax Importer understands. Importer is a simple CLI tool, and these markers are the key to make all the import and export to happen. There are several types of markers.

| Name                 | Description                                                    |
| -------------------- | -------------------------------------------------------------- |
| Importer Marker      | Main marker, used to import data from other file.              |
| Exporter Marker      | Supplemental marker used to define line range in target files. |
| Skip Importer Update | Special marker to suppress `importer update`.                  |
| Auto Generated Note  | Special marker for `importer generate` information.            |

<!-- == imptr: basic-marker / end == -->

You can find more about the markers [here](/docs/details/markers.md).

## 🚀 Examples

### `importer preview`

<!-- == imptr: example-preview / begin from: ./docs/getting-started/examples-yaml.md#[preview] == -->

`importer preview` command gives you a quick look at how the file may change when `importer update` and `importer purge` are run against the provided file. This is meant to be useful for testing and debugging.

```console
$ importer preview ./testdata/yaml/demo-before.yaml
---------------------------------------
Content Before:
1:      title: Demo of YAML Importer
2:      # == imptr: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      dummy: This will be replaced
4:      # == imptr: description / end ==
---------------------------------------

---------------------------------------
Content After Purged:
1:      title: Demo of YAML Importer
2:      # == imptr: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      # == imptr: description / end ==
---------------------------------------

---------------------------------------
Content After Processed:
1:      title: Demo of YAML Importer
2:      # == imptr: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      description: |
4:        This demonstrates how importing YAML snippet is made possible, without
5:        changing YAML handling at all.
6:      # == imptr: description / end ==
---------------------------------------

You can replace the file content with either of the commands below:

  importer update ./testdata/yaml/demo-before.yaml     Replace the file content with the Importer processed file.
  importer purge ./testdata/yaml/demo-before.yaml      Replace the file content by removing all data between marker pairs.

You can find more with 'importer help'
```

<!-- == imptr: example-preview / end == -->

### `importer generate`

<!-- == imptr: example-generate / begin from: ./docs/getting-started/examples-yaml.md#[generate] == -->

![Generate in Action][generate-in-action]

[generate-in-action]: /assets/images/importer-generate-yaml-demo.gif "Generate in Action"

<!-- == imptr: example-generate / end == -->

### `importer purge`

<!-- == imptr: example-purge / begin from: ./docs/getting-started/examples-yaml.md#[purge] == -->

![Purge in Action][purge-in-action]

[purge-in-action]: /assets/images/importer-purge-yaml-demo.mp4 "Purge in Action"

<!-- == imptr: example-purge / end == -->

You can find more examples:

- [For Markdown](/docs/getting-started/examples-markdown.md)
- [For YAML](/docs/getting-started/examples-yaml.md)

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
        run: importer update README.md
      - name: Check if README.md has any change compared to the branch
        run: |
          git status --short
          git diff-index --quiet HEAD
```

This repository uses Importer to generate some of the markdown documentation.

You can find actually running CI setup in [`.github/workflows/importer-ci.yaml`](/.github/workflows/importer-ci.yaml).

<!-- == imptr: getting-started-github-action / end == -->

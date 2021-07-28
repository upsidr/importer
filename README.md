# Importer

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

## 🚀 Examples

<!-- == imptr: getting-started-example-short / begin from: ./docs/getting-started/examples-yaml.md#[preview] == -->

```console
$ importer preview ./testdata/yaml/demo-before.yaml
---------------------------------------
Content Before:
1:      title: Demo of YAML Importer
2:      # == imptr: description / begin from: ./description-snippet.yaml#[for-demo] ==
3:      dummy: This will be replaced
4:      # == imptr: description / end ==
---------------------------------------

---------------------------------------
Content After Purged:
1:      title: Demo of YAML Importer
2:      # == imptr: description / begin from: ./description-snippet.yaml#[for-demo] ==
3:      # == imptr: description / end ==
---------------------------------------

---------------------------------------
Content After Processed:
1:      title: Demo of YAML Importer
2:      # == imptr: description / begin from: ./description-snippet.yaml#[for-demo] ==
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

<!-- == imptr: getting-started-example-short / end == -->

You can find more examples [here](https://github.com/upsidr/importer/blob/main/docs/getting-started/examples-markdown.md).

## 🧩 Supported Files

<!-- == imptr: supported-files / begin from: ./docs/details/supported-files.md#[list] == -->

| File Type | Is Supported? | File Extensions | Additional Importer Option |
| --------- | :-----------: | --------------- | -------------------------- |
| Markdown  |      ✅       | `.md`           |                            |
| YAML      |      ✅       | `.yaml`, `.yml` | Indentation                |
| HTML      |      🚧       | TBC             |                            |
| TOML      |      🚧       | TBC             |                            |

Any other file type not specified above are not supported by Importer at the moment.

For requesting additional support, [please file an issue from here](https://github.com/upsidr/importer/issues/new?assignees=&labels=enhancement&template=feature-request.yaml&title=%5BFeature+Request%5D%3A+).

<!-- == imptr: supported-files / end == -->

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

You can find actually running CI setup in [`.github/workflows/importer-markdown-ci.yaml`](https://github.com/upsidr/importer/blob/main/.github/workflows/importer-markdown-ci.yaml).

<!-- == imptr: getting-started-github-action / end == -->

## 🖋 Markers

### Importer Marker

<!-- == imptr: basic-marker / begin from: ./docs/getting-started/markers.md#[basic-marker] == -->

A marker is a simple comment with special syntax, and thus is slightly different depending on file used.

The below is a simple example for **Markdown**.

```markdown
<!-- == imptr: getting-started-install / begin from: ./docs/getting-started/install.md#[homebrew-install] == -->
```

![Marker explained][marker-explanation]

[marker-explanation]: /assets/images/marker-explanation.png "Marker Explanation"

And there has to be a matching "end" marker. This is much simpler, as options are all defined in the "begin" marker.

```markdown
<!-- == imptr: getting-started-install / end == -->
```

<!-- == imptr: basic-marker / end == -->

You can find more about the Importer Marker [here](./docs/getting-started/markers.md).

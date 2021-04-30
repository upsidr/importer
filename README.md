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

<!-- == imptr: getting-started-example-short / begin from: ./docs/getting-started/examples.md#[simple] == -->

Let's see what Importer does with the file in this repository [`./testdata/simple-before.md`](https://raw.githubusercontent.com/upsidr/importer/main/testdata/simple-before.md).

```markdown
# Simple Markdown Test

<!-- == imptr: lorem / begin from: ../docs/template/_lorem.md#5~12 == -->

Any content here will be removed by Importer.

<!-- == imptr: lorem / end == -->

Content after annotation is left untouched.
```

When you run `importer purge ./testdata/simple-before.md`:

```bash
$ importer purge ./testdata/simple-before.md
$ cat ./testdata/simple-before.md
# Simple Markdown Test

<!-- == imptr: lorem / begin from: ../docs/template/_lorem.md#5~12 == -->
<!-- == imptr: lorem / end == -->

Content after annotation is left untouched.
```

When you run `importer generate ./testdata/simple-before.md`:

```bash
$ importer generate ./testdata/simple-before.md
$ cat ./testdata/simple-before.md
# Simple Markdown Test

<!-- == imptr: lorem / begin from: ../docs/template/_lorem.md#5~12 == -->
"Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident,
sunt in culpa qui officia deserunt mollit anim id est laborum."
<!-- == imptr: lorem / end == -->

Content after annotation is left untouched.
```

<!-- == imptr: getting-started-example-short / end == -->

You can find more examples [here](https://github.com/upsidr/importer/blob/main/docs/getting-started/examples.md).

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

## 🖋 Annotations and Markers

### Importer Annotation

<!-- == imptr: basic-annotation / begin from: ./docs/getting-started/annotations.md#[basic-annotation] == -->

An annotation is a simple comment with special syntax, and thus is slightly different based on the file.

The below is a simple example for **Markdown**.

```markdown
<!-- == imptr: getting-started-install / begin from: ./docs/getting-started/install.md#[homebrew-install] == -->
```

![Annotation explained][annotation-explanation]

[annotation-explanation]: /assets/images/annotation-explanation.png "Annotation Explanation"

And there has to be a matching "end" annotation. This is much simpler, as options are all defined in the "begin" annotation.

```markdown
<!-- == imptr: getting-started-install / end == -->
```

<!-- == imptr: basic-annotation / end == -->

[You can find more about the Importer Annotation here.](./docs/getting-started/annotations.md)

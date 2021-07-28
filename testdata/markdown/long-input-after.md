# Importer

<!-- == imptr: getting-started-install / begin from: ../../docs/getting-started/install.md#[homebrew-install] == -->

You can get Importer with simple Homebrew command.

```bash
$ brew install upsidr/tap/importer
```

You can also find the relevent binary files under [releases](https://github.com/upsidr/importer/releases).

<!-- == imptr: getting-started-install / end == -->

<!-- == imptr: getting-started-example-short / begin from: ../../docs/getting-started/examples-markdown.md#[steps] == -->

**COMMAND**: Check file content before processing

```bash
cat ./testdata/markdown/demo-before.md
```

```markdown
# Markdown Demo

<!-- == imptr: short-description / begin from: ./description-snippet.md#[for-demo] == -->
Any content here will be replaced by Importer.
<!-- == imptr: short-description / end == -->
```

**COMMAND**: Preview how Importer processes the above file

```bash
importer preview ./testdata/markdown/demo-before.md
```

```console
---------------------------------------
Content Before:
1:      # Markdown Demo
2:
3:      <!-- == imptr: short-description / begin from: ./description-snippet.md#[for-demo] == -->
4:      Any content here will be replaced by Importer.
5:      <!-- == imptr: short-description / end == -->
---------------------------------------

---------------------------------------
Content After Purged:
1:      # Markdown Demo
2:
3:      <!-- == imptr: short-description / begin from: ./description-snippet.md#[for-demo] == -->
4:      <!-- == imptr: short-description / end == -->
---------------------------------------

---------------------------------------
Content After Processed:
1:      # Markdown Demo
2:
3:      <!-- == imptr: short-description / begin from: ./description-snippet.md#[for-demo] == -->
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

**COMMAND**: Update file with Importer processing

```bash
{
  cp ./testdata/markdown/demo-before.md ./testdata/markdown/demo-updated.md
  importer update ./testdata/markdown/demo-updated.md
  cat ./testdata/markdown/demo-updated.md
}
```

```markdown
# Markdown Demo

<!-- == imptr: short-description / begin from: ./description-snippet.md#[for-demo] == -->
This demonstrates how a markdown can import other file content.

Importer is a CLI tool to read and process Importer and Exporter markers.  
This can be easily integrated into CI/CD and automation setup.
<!-- == imptr: short-description / end == -->
```

You can find this file [`./testdata/markdown/demo-before.md`](https://raw.githubusercontent.com/upsidr/importer/main/testdata/markdown/demo-before.md).

<!-- == imptr: getting-started-example-short / end == -->

You can find more examples [here](https://github.com/upsidr/importer/blob/main/docs/getting-started/examples-markdown.md).

<!-- == imptr: getting-started-github-action / begin from: ../../docs/getting-started/github-actions.md#1~32 == -->
## :octocat: GitHub Action Integration


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

<!-- == imptr: some_random_note / begin from: ../../docs/template/_lorem.md#5~12 == -->
"Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident,
sunt in culpa qui officia deserunt mollit anim id est laborum."
<!-- == imptr: some_random_note / end == -->

<!-- == imptr: import_from_proposal / begin from: ../../Proposal.md#5~8 == -->
## What it does

- In any file, you can import other file content
- Works like code generation, but does not need to be code
<!-- == imptr: import_from_proposal / end == -->

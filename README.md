# Importer

<!-- == imptr: getting-started-install / begin from: ./docs/getting-started/install.md#1~9 == -->
## ✨ Install

The simplest approach is to get Importer via Homebrew.

```bash
$ brew install upsidr/tap/importer
```

You can also find the relevent binary files under [releases](https://github.com/upsidr/importer/releases).
<!-- == imptr: getting-started-install / end == -->

<!-- Temporarily disabling Example section -->
<!-- == DISABLED imptr: getting-started-example-short / begin from: ./docs/getting-started/examples.md#1~49 == -->

## 🚀 Examples

<!-- == DISABLED imptr: getting-started-example-short / end == -->

You can find more examples [here](https://github.com/upsidr/importer/blob/main/docs/getting-started/examples.md).

<!-- == imptr: getting-started-github-action / begin from: ./docs/getting-started/github-actions.md#1~30 == -->
## :octocat: GitHub Action Integration

Because you can install Importer using Homebrew, you can set up GitHub Action definition such as below:

<!--TODO: This is exactly where Importer can pull in the actual file content-->

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

<!-- == imptr: import_from_proposal / begin from: ./Proposal.md#5~8 == -->
## What it does

- In any file, you can import other file content
- Works like code generation, but does not need to be code
<!-- == imptr: import_from_proposal / end == -->

<!-- == imptr: some_random_note / begin from: ./docs/template/_lorem.md#5~12 == -->
"Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident,
sunt in culpa qui officia deserunt mollit anim id est laborum."
<!-- == imptr: some_random_note / end == -->

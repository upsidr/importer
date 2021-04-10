# Importer

<!-- == imptr: getting-started-install / begin from: ./docs/getting-started/install.md#1~9 == -->
## ✨ Install

The simplest approach is to get Importer via Homebrew.
 
```bash
$ brew install upsidr/tap/importer
```

You can also find the relevent binary files under [releases](https://github.com/upsidr/importer/releases).
<!-- == imptr: getting-started-install / end == -->

<!-- == imptr: getting-started-example-short / begin from: ./docs/getting-started/examples.md#1~49 == -->
## 🚀 Examples

Let's see what Importer does with the file in this repository [`./testdata/simple-before.md`](https://raw.githubusercontent.com/upsidr/importer/main/testdata/simple-before.md).

```console
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

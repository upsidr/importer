# importer

<!-- == imptr: getting-started / begin from: ./docs/getting-started.md#3~63 == -->
## Basic Usage

```bash
$ cat testdata/simple-before-importer.md
```

```console
# Simple Markdown Test

<!-- == imptr: lorem / begin from: ./docs/template/_lorem.md#5~12 == -->

Any content here will be removed by Importer.

<!-- == imptr: lorem / end == -->
```

```bash
$ importer preview testdata/simple-before-importer.md
```

```console
---------------
Content Before:
0:      # Simple Markdown Test
1:
2:      <!-- == imptr: lorem / begin from: ./docs/template/_lorem.md#5~12 == -->
3:
4:      Any content here will be removed by Importer.
5:
6:      <!-- == imptr: lorem / end == -->
---------------

---------------
Content After Purged:
0:      # Simple Markdown Test
1:
2:      <!-- == imptr: lorem / begin from: ./docs/template/_lorem.md#5~12 == -->
3:      <!-- == imptr: lorem / end == -->
---------------

---------------
Content After Processed:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ./docs/template/_lorem.md#5~12 == -->
4:      "Lorem ipsum dolor sit amet,
5:      consectetur adipiscing elit,
6:      sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
7:      Ut enim ad minim veniam,
8:      quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
9:      Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
10:     Excepteur sint occaecat cupidatat non proident,
11:     sunt in culpa qui officia deserunt mollit anim id est laborum."
12:     <!-- == imptr: lorem / end == -->
---------------

You can replace the file content with either of the commands below:

- 'importer generate testdata/simple-before-importer.md'
- 'importer purge testdata/simple-before-importer.md'
```
<!-- == imptr: getting-started / end == -->

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

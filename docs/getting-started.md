# Getting Started with Importer

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
0:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ./docs/template/_lorem.md#5~12 == -->
4:
5:      Any content here will be removed by Importer.
6:
7:      <!-- == imptr: lorem / end == -->
---------------

---------------
Content After Purged:
0:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ./docs/template/_lorem.md#5~12 == -->
4:      <!-- == imptr: lorem / end == -->
---------------

---------------
Content After Processed:
1:
2:      # Simple Markdown Test
3:
4:      <!-- == imptr: lorem / begin from: ./docs/template/_lorem.md#5~12 == -->
5:      "Lorem ipsum dolor sit amet,
6:      consectetur adipiscing elit,
7:      sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
8:      Ut enim ad minim veniam,
9:      quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
10:     Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
11:     Excepteur sint occaecat cupidatat non proident,
12:     sunt in culpa qui officia deserunt mollit anim id est laborum."
13:     <!-- == imptr: lorem / end == -->
---------------
```

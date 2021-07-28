## 🚀 Examples

### With Markdown

<!-- == export: simple-markdown / begin == -->

**COMMAND**: Check file content before processing

```bash
cat ./testdata/markdown/demo-before.md
```

**OUTPUT**

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

**OUTPUT**

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

**OUTPUT**

```markdown
# Markdown Demo

<!-- == imptr: short-description / begin from: ./description-snippet.md#[for-demo] == -->
This demonstrates how a markdown can import other file content.

Importer is a CLI tool to read and process Importer and Exporter markers.  
This can be easily integrated into CI/CD and automation setup.
<!-- == imptr: short-description / end == -->
```

You can find this file [`./testdata/markdown/demo-before.md`](https://raw.githubusercontent.com/upsidr/importer/main/testdata/markdown/demo-before.md).

<!-- == export: simple-markdown / end == -->

### Full Example

The below allows you to experiment Importer offering without cloning this repository.

```bash
{
    # Create the above file in temp directory
    cat << EOF > /tmp/importer-example.md
# Simple Markdown Test

<!-- == imptr: lorem / begin from: ./lorem.md#5~12 == -->

Any content here will be removed by Importer.

<!-- == imptr: lorem / end == -->

Content after marker is left untouched.
EOF

    # Create a file with Lorem Ipsum in a separate file
    cat << EOF > /tmp/lorem.md
# Test Note

This file contains note that's used in other markdown files.

"Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident,
sunt in culpa qui officia deserunt mollit anim id est laborum."
EOF
}
```

Importer currently supports 3 commands:

- `importer preview`
- `importer purge`
- `importer generate`

Preview allows you to see how Importer processed the file.

```bash
$ importer preview /tmp/importer-example.md
```

<details>

<summary>Expand to see the full output</summary>

```console
---------------
Content Before:
0:      # Simple Markdown Test
1:
2:      <!-- == imptr: lorem / begin from: ../../docs/template/_lorem.md#5~12 == -->
3:
4:      Any content here will be removed by Importer.
5:
6:      <!-- == imptr: lorem / end == -->
7:
8:      Content after marker is left untouched.
---------------

---------------
Content After Purged:
0:      # Simple Markdown Test
1:
2:      <!-- == imptr: lorem / begin from: ../../docs/template/_lorem.md#5~12 == -->
3:      <!-- == imptr: lorem / end == -->
4:
5:      Content after marker is left untouched.
---------------

---------------
Content After Processed:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ../../docs/template/_lorem.md#5~12 == -->
4:      "Lorem ipsum dolor sit amet,
5:      consectetur adipiscing elit,
6:      sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
7:      Ut enim ad minim veniam,
8:      quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
9:      Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
10:     Excepteur sint occaecat cupidatat non proident,
11:     sunt in culpa qui officia deserunt mollit anim id est laborum."
12:     <!-- == imptr: lorem / end == -->
13:
14:     Content after marker is left untouched.
---------------

You can replace the file content with either of the commands below:

- 'importer generate testdata/simple-before.md'
  Replace the file content with the processed file, importing all annotated references.
- 'importer purge testdata/simple-before.md'
  Replace the file content by removing all data between marker pairs.

You can find more with 'importer help'
```

</details>

With the `importer preview` command, you get the idea of how the file is going to look like.

The below is how the file would look like after `importer purge` and `importer generate`.

#### `importer purge`

```bash
{
    importer purge /tmp/importer-example.md
    cat /tmp/importer-example.md
}
```

<details>

<summary>Expand to see the full output</summary>

```console
# Simple Markdown Test

<!-- == imptr: lorem / begin from: ./lorem.md#5~12 == -->
<!-- == imptr: lorem / end == -->

Content after marker is left untouched.
```

</details>

#### `importer generate`

<details>

<summary>Expand to see the full output</summary>

```console
cat /tmp/importer-example.md
# Simple Markdown Test

<!-- == imptr: lorem / begin from: ./lorem.md#5~12 == -->
"Lorem ipsum dolor sit amet,
consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam,
quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.
Excepteur sint occaecat cupidatat non proident,
sunt in culpa qui officia deserunt mollit anim id est laborum."
<!-- == imptr: lorem / end == -->

Content after marker is left untouched.
```

</details>

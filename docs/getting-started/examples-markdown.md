## 🚀 Examples with Markdown

### Preview

<!-- == export: preview / begin == -->

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

<!-- == export: preview / end == -->

### Update

<!-- == export: update / begin == -->

`importer update` imports based on Importer Markers in the given file, and update the file in place. This is useful for having a single file to manage and also import other file contents. If you want to have a template file which only holds Importer Markers and not actually the imported content, you should use `importer generate` instead.

> 🕹 COMMAND

```bash
# Check demo file before update
cat ./testdata/markdown/demo-before.md
```

```markdown
# Markdown Demo

<!-- == imptr: short-description / begin from: ./snippet-description.md#[for-demo] == -->

Any content here will be replaced by Importer.

<!-- == imptr: short-description / end == -->
```

> 🕹 COMMAND

```bash
# Update file with Importer processing.
#
# Because Importer updates the file in place, this is making a copy of the
# "-before" file, and running importer update against the copied file of
# "-updated" file.
{
  cp ./testdata/markdown/demo-before.md ./testdata/markdown/demo-updated.md
  importer update ./testdata/markdown/demo-updated.md
  cat ./testdata/markdown/demo-updated.md
}
```

```markdown
# Markdown Demo

<!-- == imptr: short-description / begin from: ./snippet-description.md#[for-demo] == -->

This demonstrates how a markdown can import other file content.

Importer is a CLI tool to read and process Importer and Exporter markers.  
This can be easily integrated into CI/CD and automation setup.

<!-- == imptr: short-description / end == -->
```

You can find these files here:

- [`/testdata/markdown/demo-before.md`](https://raw.githubusercontent.com/upsidr/importer/main/testdata/markdown/demo-before.md)
- [`/testdata/markdown/demo-updated.md`](https://raw.githubusercontent.com/upsidr/importer/main/testdata/markdown/demo-updated.md)

<!-- == export: update / end == -->

### Purge

<!-- == export: purge / begin == -->

`importer purge` removes any lines between Importer Markers in the given file, and update the file in place. The same operation is executed for `importer update` before importing all the lines, but this "purge" is sometimes useful to see the file without extra data imported.

> 🕹 COMMAND

```bash
# Check demo file before update
cat ./testdata/markdown/demo-before.md
```

```markdown
# Markdown Demo

<!-- == imptr: short-description / begin from: ./snippet-description.md#[for-demo] == -->

Any content here will be replaced by Importer.

<!-- == imptr: short-description / end == -->
```

> 🕹 COMMAND

```bash
# Purge any text between Importer Markers.
#
# Because Importer updates the file in place, this is making a copy of the
# "-before" file, and running importer update against the copied file of
# "-puged" file.
{
  cp ./testdata/markdown/demo-before.md ./testdata/markdown/demo-purged.md
  importer purge ./testdata/markdown/demo-purged.md
  cat ./testdata/markdown/demo-purged.md
}
```

```markdown
# Markdown Demo

<!-- == imptr: short-description / begin from: ./snippet-description.md#[for-demo] == -->
<!-- == imptr: short-description / end == -->
```

You can find these files here:

- [`/testdata/markdown/demo-before.md`](https://raw.githubusercontent.com/upsidr/importer/main/testdata/markdown/demo-before.md)
- [`/testdata/markdown/demo-purged.md`](https://raw.githubusercontent.com/upsidr/importer/main/testdata/markdown/demo-purged.md)

<!-- == export: update / end == -->

### Generate

<!-- == export: generate / end == -->

`importer generate` imports based on Importer Markers in the given file, and write the result to stdout or file. This can be used for debugging, or create a template file with Importer Markers but keep the file purely for Importer Markers.

> 🕹 COMMAND

```bash
# Check demo file before update
cat ./testdata/markdown/demo-before.md
```

```markdown
# Markdown Demo

<!-- == imptr: short-description / begin from: ./snippet-description.md#[for-demo] == -->

Any content here will be replaced by Importer.

<!-- == imptr: short-description / end == -->
```

> 🕹 COMMAND

```bash
# Check the result of Importer processing
importer generate ./testdata/markdown/demo-before.md

# If you want to write to a file, you can provide --out FILENAME.
# TODO: the below command doesn't work due to the flag handling, --out needs to be before the filename for it to work as of v0.0.1-rc7
# importer generate ./testdata/markdown/demo-before.md --out some-target.md
```

```markdown
# Markdown Demo

<!-- == imptr: short-description / begin from: ./snippet-description.md#[for-demo] == -->

This demonstrates how a markdown can import other file content.

Importer is a CLI tool to read and process Importer and Exporter markers.  
This can be easily integrated into CI/CD and automation setup.

<!-- == imptr: short-description / end == -->
```

You can find this files here:

- [`/testdata/markdown/demo-before.md`](https://raw.githubusercontent.com/upsidr/importer/main/testdata/markdown/demo-before.md)

<!-- == export: generate / end == -->

## 🎯 Full Example

The below steps allow you to experiment Importer offerings without cloning this repository.

You need to have Importer installed.

```bash
{
    # Create the above file in temp directory
    cat << EOF > /tmp/importer-example.md
# Simple Markdown Test

<!-- == imptr: lorem / begin from: ./snippet-lorem.md#5~12 == -->

Any content here will be removed by Importer.

<!-- == imptr: lorem / end == -->

Content after marker is left untouched.
EOF

    # Create a file with Lorem Ipsum in a separate file
    cat << EOF > /tmp/snippet-lorem.md
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

Let's check out 3 commands below.

- `importer preview`
- `importer purge`
- `importer update`

#### `importer preview`

Preview allows you to see how Importer processed the file.

```bash
importer preview /tmp/importer-example.md
```

<details>

<summary>Expand to see the full output</summary>

```console
---------------------------------------
Content Before:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ./snippet-lorem.md#5~12 == -->
4:
5:      Any content here will be removed by Importer.
6:
7:      <!-- == imptr: lorem / end == -->
8:
9:      Content after marker is left untouched.
---------------------------------------

---------------------------------------
Content After Purged:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ./snippet-lorem.md#5~12 == -->
4:      <!-- == imptr: lorem / end == -->
5:
6:      Content after marker is left untouched.
---------------------------------------

---------------------------------------
Content After Processed:
1:      # Simple Markdown Test
2:
3:      <!-- == imptr: lorem / begin from: ./snippet-lorem.md#5~12 == -->
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
---------------------------------------

You can replace the file content with either of the commands below:

  importer update /tmp/importer-example.md     Replace the file content with the Importer processed file.
  importer purge /tmp/importer-example.md      Replace the file content by removing all data between marker pairs.

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

<!-- == imptr: lorem / begin from: ./snippet-lorem.md#5~12 == -->
<!-- == imptr: lorem / end == -->

Content after marker is left untouched.
```

</details>

#### `importer update`

```bash
{
    importer update /tmp/importer-example.md
    cat /tmp/importer-example.md
}
```

<details>

<summary>Expand to see the full output</summary>

```console
cat /tmp/importer-example.md
# Simple Markdown Test

<!-- == imptr: lorem / begin from: ./snippet-lorem.md#5~12 == -->
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

<!-- == importer-skip-update == -->

# Markers

**Markers** are a simple comment with special syntax Importer understands. There
are several types of markers.

## Basic Syntax

The markers always follow the pattern of `== some-importer-marker-input ==`.

In case of YAML, this would be `# == some-importer-marker-input ==`.\
In case of markdown, this would be `<!-- == some-importer-marker-input == -->`.

The main markers **Importer Markers** and **Exporter Markers** are both made up of pairs, `begin` and `end`.

Other markers are used to update Importer behaviours.

## Importer Marker

![Importer Marker Syntax](/assets/images/importer-marker-syntax.png)

> NOTE: The above example is from [`/testdata/markdown/simple-before.md`](/testdata/markdown/simple-before.md).

### 1️⃣ Importer Marker Type

- Tell Importer to start the import setup.
- This can be represented with `importer`, `import`, `imptr` or `i`.
- Do not forget to add `:` at the end.

### 2️⃣ Importer Marker Name

- Any name of your choice, with no whitespace character.
- The same name cannot be used in a single file.

### ℹ️ Separate with ` / `

- Add separator using ` / `. The spaces around the `/` are required as of now.

### 3️⃣ Either `begin` or `end`

- Each Importer Marker must be a pair to operate.

### 4️⃣ Importer Marker Details

This includes target file to import from, etc.

- `from: FILENAME#OPTION`: Define where import from.
  - `FILENAME`: Import from the `FILENAME`, the location of target file can be a URL or relative path from the source file.
  - `OPTION`: Define which lines to import.
    - `NUM1~NUM2`: Import line range from `NUM1` to `NUM2`. Leaving `NUM1` empty means from the beginning of the file. Leaving `NUM2` empty means to the end of the file.
    - `NUM1,NUM2`: Import each lines specified (e.g. `NUM1`, `NUM2`) one by one.
    - `[Exporter-Marker]`: Import lines based on Exporter Markers defined in the target file.
- `indent: [align|absolute NUM|extra NUM|keep]`: Update indentation for the imported data.
  - `align`: Align to the indentation of Importer Marker.
  - `absolute NUM` (e.g. `absolute 2`): Update indentation to `NUM` spaces. This ignores the original indentation from the imported data, but keeps the tree structure.
  - `extra NUM` (e.g. `extra 4`): Add extra indentation of `NUM` spaces.
  - `keep` (default): Keep the indentation from the imported data.

### Examples

#### With `/testdata/markdown/simple-before.md`

<details>
<summary>Preview Importer CLI in action</summary>

```console
$ importer preview testdata/markdown/simple-before.md
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

  importer update testdata/markdown/simple-before.md     Replace the file content with the Importer processed file.
  importer purge testdata/markdown/simple-before.md      Replace the file content by removing all data between marker pairs.

You can find more with 'importer help'
```

</details>

#### With `/testdata/yaml/snippet-description.yaml` using Exporter Marker

<details>
<summary>Preview Importer CLI in action</summary>

```console
$ cat testdata/yaml/snippet-description.yaml
# == export: for-demo / begin ==
description: |
  This demonstrates how importing YAML snippet is made possible, without
  changing YAML handling at all.
# == export: for-demo / end ==

$ importer preview testdata/yaml/demo-before.yaml
---------------------------------------
Content Before:
1:      title: Demo of YAML Importer
2:      # == import: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      dummy: This will be replaced
4:      # == import: description / end ==
---------------------------------------

---------------------------------------
Content After Purged:
1:      title: Demo of YAML Importer
2:      # == import: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      # == import: description / end ==
---------------------------------------

---------------------------------------
Content After Processed:
1:      title: Demo of YAML Importer
2:      # == import: description / begin from: ./snippet-description.yaml#[for-demo] ==
3:      description: |
4:        This demonstrates how importing YAML snippet is made possible, without
5:        changing YAML handling at all.
6:      # == import: description / end ==
---------------------------------------

You can replace the file content with either of the commands below:

  importer update testdata/yaml/demo-before.yaml     Replace the file content with the Importer processed file.
  importer purge testdata/yaml/demo-before.yaml      Replace the file content by removing all data between marker pairs.

You can find more with 'importer help'
```

</details>

---

### Exporter Marker

![Exporter Marker Syntax](/assets/images/exporter-marker-syntax.png)

> NOTE: The above example is from [`/testdata/yaml/snippet-k8s-resource.yaml`](/testdata/yaml/snippet-k8s-resource.yaml).

#### 1️⃣ Exporter Marker Type

- Tell Importer to start the export setup.
- This can be represented with `exporter`, `export`, `exptr` or `e`.

#### 2️⃣ Importer Marker Name

- Any name of your choice, with no whitespace character.
- The same name cannot be used in a single file.

#### 3️⃣ Either `begin` or `end`

- Each Exporter Marker must be a pair to operate.

---

The main marker for importing data from other file.

This needs to be closed with `== import: NAME / end ==`.

### Exporter Marker: `== export: NAME / begin ==`

Exporter Markers can be used to mark specific lines as import target. This allows Importer to not specify the line range, but simply rely on the exporter markers to find which lines to import.

### List

The marker pairs are the most common and useful markers.

The pair needs to have the same "name", and should have "begin" and "end".

| Name                    | Syntax                                      | Use Case                              |
| ----------------------- | ------------------------------------------- | ------------------------------------- |
| Importer Marker (begin) | `== import: NAME / begin from: FILENAME ==` | Import data from `FILENAME`.          |
| Importer Marker (end)   | `== import: NAME / end ==`                  | Close Importer Marker.                |
| Exporter Marker (begin) | `== export: NAME / begin ==`                | Mark specific lines as import target. |
| Exporter Marker (end)   | `== export: NAME / end ==`                  | Close Exporter Marker.                |

## Marker Types - Special Markers

There are some special markers that are used to update Importer behaviour.

| Name                 | Syntax                                    | Use Case                                                                                                                                       |
| -------------------- | ----------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| Auto Generated Note  | `== improter-generated-from: FILENAME ==` | This is auto-generated by `importer generate FILENAME --out TARGET_FILE`, which tells how the file was generated by using `FILENAME` as input. |
| Skip Importer Update | `== importer-skip-update ==`              | Mark the file not to be updated by `importer update` command.                                                                                  |

## 📌 Importer Marker

Importer's most important syntax is its markers. Each marker has `begin` and
`end` to ensure easy usage and reproducibility.

### Basic Syntax

<!-- == export: basic-marker / begin == -->

A marker is a simple comment with special syntax, and thus is slightly different
depending on file used.

The below is a simple example for **Markdown**.

```markdown
<!-- == imptr: getting-started-install / begin from: ./docs/getting-started/install.md#[homebrew-install] == -->
```

![Marker explained][marker-explanation]

[marker-explanation]: /assets/images/marker-explanation.png "Marker Explanation"

And there has to be a matching "end" marker. This is much simpler, as options
are all defined in the "begin" marker.

```markdown
<!-- == imptr: getting-started-install / end == -->
```

<!-- == export: basic-marker / end == -->

### Importer Syntax in Detail

As described above, Importer syntax has 3 parts.

- Importer Name
- Importer "begin" or "end"
- Importer Options

#### Importer Name

> Above Example: `getting-started-install`

Importer Name can be any string, without any whitespace characters.

#### Importer "begin" or "end"

There are only 2 options, "begin" or "end". It has to be lower case.

#### Importer Options

> Above Example: `from: ./docs/getting-started/install.md#[homebrew-install]`

Importer Options are anything you define between `/` and `==`.

The below are the Option formats:

| Name                       | Example         | Description                                                                                                                                                                                                                                                                                                 |
| -------------------------- | --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Target Path                | `from: xyz.md`  | Defines where to import from. This is a relative path from the file containing the marker.<br /><br /> **Known Limitations**: Path cannot contain whitespace characters.                                                                                                                                    |
| Separator                  | `#`             | This is to separate Target Path and Target Detail. It can have as many preceding whispace characters.                                                                                                                                                                                                       |
| Target Detail - Line Range | `[1~33]`        | Imports only provided line ranges. You can omit before or after `~` to indicate the range starts from the beginning of the file, or ends at the end of the file.                                                                                                                                            |
| Target Detail - Line List  | `[1,2,5]`       | Imports only provided lines. The lines are comma separated, and you can also use line range in the same target detail. <br /><br /> **Known Limitations**: The order of lines is not persisted, and thus if you define `[3,2,1]`, you would actually see lines imported as line#1, line#2, and then line#3. |
| Target Detail - Marker     | `[some-marker]` | Searches for the matching Export Marker in the target file. More about Export Marke below. <br /><br /> **Known Limitations**: You can only provide single marker.                                                                                                                                          |

### Exporter Marker

Importer's simplest form is to import some lines from another file by providing
line numbers.

But if you want to import file that gets updated frequently, it is quite
cumbersome to make line number adjustments on all the Importer Markers every
time.

Exporter Marker allows defining the begin / end for line range, and assigns a
name to it.

```markdown
<!-- == export: some-export-name / begin == -->

This is the data that can be exported under the Export name of `some-export-name`.

<!-- == export: some-export-name / end == -->

Any content before or after the marker is not imported.
```

With the above Exporter Marker, you can then use something like below to import:

```markdown
<!-- == imptr: export-marker-test / begin from: ./target-file.md#[some-export-name] == -->
<!-- == imptr: export-marker-test / end == -->
```

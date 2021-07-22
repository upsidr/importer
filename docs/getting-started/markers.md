<!-- == imptr: ignore-all == -->

# Markers

## 📌 Importer Marker

Importer's most important syntax is its markers. Each marker has `begin` and `end` to ensure easy usage and reproducibility.

### Basic Syntax

<!-- == export: basic-marker / begin == -->

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

Importer's simplest form is to import some lines from another file by providing line numbers.

But if you want to import file that gets updated frequently, it is quite cumbersome to make line number adjustments on all the Importer Markers every time.

Exporter Marker allows defining the begin / end for line range, and assigns a name to it.

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

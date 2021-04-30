<!-- == imptr: ignore-all == -->

# Annotations

## 📌 Importer Annotation

Importer's most important syntax is its annotation. Each annotation has `begin` and `end` to ensure easy usage and reproducibility.

### Basic Syntax

<!-- == export: basic-annotation / begin == -->

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

<!-- == export: basic-annotation / end == -->

### Importer Syntax in Detail

As described above, Importer syntax has 3 parts.

- Importer Name
- Importer "begin" or "end"
- Importer Options

#### Importer Name

Importer Name can be any string, without any whitespace characters.

#### Importer "begin" or "end"

There are only 2 options, "begin" or "end". It has to be lower case.

#### Importer Options

Importer Options are anything you define between `/` and `==`. More about the options below.

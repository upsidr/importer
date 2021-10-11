## 🎮 Commands

<!-- == export: help-output / begin == -->

```console
$ importer -h
Import any lines, from anywhere

Usage:
  importer [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  generate    Processes Importer markers and send output to stdout or file
  help        Help about any command
  preview     Shows a preview of Importer update and purge results
  purge       Removes all imported lines and update the file in place
  update      Processes Importer markers and update the file in place

Flags:
  -h, --help   help for importer

Use "importer [command] --help" for more information about a command.
```

<!-- == export: help-output / end == -->

<!-- == export: list / begin == -->

| Name                     | Description                                                                                       |
| ------------------------ | ------------------------------------------------------------------------------------------------- |
| `importer preview FILE`  | Write before/purged/after preview of how Importer processes the file content to stdout.           |
| `importer update FILE`   | Run Importer processing on `FILE`, and update it in place.                                        |
| `importer purge FILE`    | Parse Importer Markers, remove any content within Importer Markers, and update the file in plcae. |
| `importer generate FILE` | Run Importer processing on `FILE`, and write the result to stdout.                                |

<!-- == export: list / end == -->

### `importer generate`

```console
$ importer generate --help

`generate` command parses the provided file as the input, and output the processed file content to stdout or a file.

While `update` command is useful for managing file content in itself, `generate` can be used to create a separate template file.
This approach allows the input file to be full of Importer markes without actual importing, and only used as the template to generate a new file.

Usage:
  importer generate [filename] [flags]

Aliases:
  generate, gen

Flags:
      --disable-header   disable automatically added header of Importer generated notice
  -h, --help             help for generate
      --keep-markers     keep Importer Markers from the generated result
  -o, --out FILE         write to FILE
```

### `importer update`

```console
$ importer update --help

`update` command parses the provided file and processes the Import markers in place.

This does not support creating a new file, nor send the result to stdout. For such use cases, use `generate` command

Usage:
  importer update [filename] [flags]

Aliases:
  update, up

Flags:
      --dry-run   Run without updating the file
  -h, --help      help for update
```

### `importer purge`

```console
$ importer purge --help

`purge` command processes the provided file and removes all the contents surrounded by Importer markers.

Importer markers will be left intact.

Usage:
  importer purge [filename] [flags]

Flags:
      --dry-run   Run without updating the file
  -h, --help      help for purge
```

### `importer preview`

```console
$ importer preview --help

`preview` command processes the provided file and gives you a quick preview.

This allows you to find what the file looks like after `update` or `purge`.

Usage:
  importer preview [filename] [flags]

Aliases:
  preview, pre, p

Flags:
  -h, --help     help for preview
      --lines    Show line numbers
  -p, --purge    Show only purged result
  -u, --update   Show only updated result
```

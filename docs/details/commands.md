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

| Name                         | Description                                                                                       |
| ---------------------------- | ------------------------------------------------------------------------------------------------- |
| `importer generate FILENAME` | Run Importer processing on `FILENAME`, and write the result to stdout.                            |
| `importer update FILENAME`   | Run Importer processing on `FILENAME`, and update it in place.                                    |
| `importer purge FILENAME`    | Parse Importer Markers, remove any content within Importer Markers, and update the file in plcae. |
| `importer preview FILENAME`  | Write before/purged/after preview of how Importer processes the file content to stdout.           |

<!-- == export: list / end == -->

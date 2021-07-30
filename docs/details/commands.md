## 🎮 Commands

<!-- == export: list / begin == -->

| Name                         | Description                                                                                       |
| ---------------------------- | ------------------------------------------------------------------------------------------------- |
| `importer generate FILENAME` | Run Importer processing on `FILENAME`, and write the result to stdout.                            |
| `importer update FILENAME`   | Run Importer processing on `FILENAME`, and update it in place.                                    |
| `importer purge FILENAME`    | Parse Importer Markers, remove any content within Importer Markers, and update the file in plcae. |
| `importer preview FILENAME`  | Write before/purged/after preview of how Importer processes the file content to stdout.           |

<!-- == export: list / end == -->


<!-- == export: help-output / begin == -->

```console
$ importer -h
NAME:
   importer - Import any lines, from anywhere

USAGE:
   importer [command]

COMMANDS:
   preview        Provides Importer update and purge previews
   update, up     Processes Importer markers and update the file in place
   generate, gen  Processes Importer markers and send output to stdout or file
   purge          Removes all imported lines and update the file in place
   help, h        Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

<!-- == export: help-output / end == -->
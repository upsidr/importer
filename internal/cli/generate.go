package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/upsidr/importer/internal/parse"
)

var (
	generateCliCmd = &cli.Command{
		Name:      "generate",
		Aliases:   []string{"gen"},
		UsageText: rootCmdName + " generate [filename]",
		Usage:     "Processes Importer markers and send output to stdout or file",
		Description: `
` + "`generate`" + ` command parses the provided file as the input, and output the processed file content to stdout or a file.

While ` + "`update`" + ` command is useful for managing file content in itself, ` + "`generate`" + ` can be used to create a separate template file.
This approach allows the input file to be full of Importer markes without actual importing, and only used as the template to generate a new file.
`,
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:        "out",
				Aliases:     []string{"o"},
				Usage:       "write to `FILE`",
				Destination: &generateTargetFile,
			},
		},
		Action: executeGenerateCLI,
	}
	generateTargetFile string
)

func executeGenerateCLI(cmd *cli.Context) error {
	args := cmd.Args()
	// TODO: add some util func to hande all common error cases
	if args.Len() < 1 {
		return errors.New("error: incorrect argument, you need to pass in an argument")
	}

	arg := args.First()
	out := generateTargetFile
	if err := generate(arg, out); err != nil {
		return fmt.Errorf("error: handling generate, %v", err)
	}

	return nil
}

func generate(fileName string, targetFilepath string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	file, err := parse.Parse(fileName, f)
	if err != nil {
		return err
	}

	err = file.ProcessAnnotations()
	if err != nil {
		return err
	}

	if targetFilepath != "" {
		return file.WriteAfterTo(targetFilepath)
	}

	return file.PrintAfter()
}

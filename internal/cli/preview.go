package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/urfave/cli/v2"

	"github.com/upsidr/importer/internal/parse"
)

var (
	previewCmd = &cobra.Command{
		Use:   "preview",
		Short: "Parse the provided file and send the result to stdout",
		RunE:  executePreview,
	}
	previewCliCmd = &cli.Command{
		Name:        "preview",
		UsageText:   rootCmdName + " preview [filename]",
		Usage:       "Parse the provided file and send the result to stdout",
		Description: "Parse the provided file and send the result to stdout",
		Action:      executePreviewCLI,
	}
)

func executePreview(cmd *cobra.Command, args []string) error {
	// TODO: add some util func to hande all common error cases
	if len(args) != 1 {
		return errors.New("error: incorrect argument, you can only pass in 1 argument")
	}

	arg := args[0]
	if err := preview(arg); err != nil {
		return fmt.Errorf("error: handling preview, %v", err)
	}

	return nil
}

func executePreviewCLI(ctx *cli.Context) error {
	args := ctx.Args()
	// TODO: add some util func to hande all common error cases
	if args.Len() != 1 {
		return errors.New("error: incorrect argument, you can only pass in 1 argument")
	}

	arg := args.First()
	if err := preview(arg); err != nil {
		return fmt.Errorf("error: handling preview, %v", err)
	}

	return nil
}

func preview(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	file, err := parse.Parse(fileName, f)
	if err != nil {
		return err
	}

	file.PrintDebugAll()

	fmt.Printf(`You can replace the file content with either of the commands below:

- 'importer generate %s'
  Replace the file content with the processed file, importing all annotated references.
- 'importer purge %s' 
  Replace the file content by removing all data between annotation pairs.

You can find more with 'importer help'
`, fileName, fileName)

	return nil
}

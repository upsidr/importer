package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/upsidr/importer/internal/parse"
)

var (
	purgeCliCmd = &cli.Command{
		Name:      "purge",
		UsageText: rootCmdName + " purge [filename]",
		Usage:     "Removes all imported lines and update the file in place",
		Description: `
` + "`purge`" + ` command processes the provided file and removes all the contents surrounded by Importer markers.

Importer markers will be left intact.
`,
		Action: executePurgeCLI,
	}
)

func executePurgeCLI(ctx *cli.Context) error {
	args := ctx.Args()
	// TODO: add some util func to hande all common error cases
	if args.Len() != 1 {
		return errors.New("error: incorrect argument, you can only pass in 1 argument")
	}

	arg := args.First()
	if err := purge(arg); err != nil {
		return fmt.Errorf("error: handling purge, %v", err)
	}

	return nil
}

func purge(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	file, err := parse.Parse(fileName, f)
	if err != nil {
		return err
	}

	err = file.ReplaceWithPurged()
	if err != nil {
		return err
	}

	return nil
}

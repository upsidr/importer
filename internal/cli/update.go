package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/upsidr/importer/internal/parse"
)

var (
	updateCmd = &cli.Command{
		Name:        "update",
		Aliases:     []string{"up"},
		UsageText:   rootCmdName + " update [filename]",
		Usage:       "Parse the provided file and generate result with imported files",
		Description: "Parse the provided file and generate result with imported files",
		Action:      executeUpdate,
	}
)

func executeUpdate(ctx *cli.Context) error {
	args := ctx.Args()
	// TODO: add some util func to hande all common error cases
	if args.Len() != 1 {
		return errors.New("error: incorrect argument, you can only pass in 1 argument")
	}

	arg := args.First()
	if err := update(arg); err != nil {
		return fmt.Errorf("error: handling generate, %v", err)
	}

	return nil
}

func update(fileName string) error {
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

	err = file.ReplaceWithAfter()
	if err != nil {
		return err
	}

	return nil
}

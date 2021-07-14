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
	purgeCmd = &cobra.Command{
		Use:   "purge",
		Short: "Parse the provided file and purge data between annotations",
		RunE:  executePurge,
	}
	purgeCliCmd = &cli.Command{
		Name:        "purge",
		UsageText:   rootCmdName + " purge [filename]",
		Usage:       "Parse the provided file and purge data between annotations",
		Description: "Parse the provided file and purge data between annotations",
		Action:      executePurgeCLI,
	}
)

func executePurge(cmd *cobra.Command, args []string) error {
	// TODO: add some util func to hande all common error cases
	if len(args) != 1 {
		return errors.New("error: incorrect argument, you can only pass in 1 argument")
	}

	arg := args[0]
	if err := purge(arg); err != nil {
		return fmt.Errorf("error: handling purge, %v", err)
	}

	return nil
}
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

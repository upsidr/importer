package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/upsidr/importer/internal/parse"
)

var (
	purgeCmd = &cobra.Command{
		Use:   "purge",
		Short: "Parse the provided file and purge data between annotations",
		RunE:  executePurge,
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

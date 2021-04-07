package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/upsidr/importer/internal/parse"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Parse the provided file and generate result with imported files",
		RunE:  executeGenerate,
	}
)

func executeGenerate(cmd *cobra.Command, args []string) error {
	// TODO: add some util func to hande all common error cases
	if len(args) != 1 {
		return errors.New("error: incorrect argument, you can only pass in 1 argument")
	}

	arg := args[0]
	if err := generate(arg); err != nil {
		return fmt.Errorf("error: handling generate, %v", err)
	}

	return nil
}

func generate(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	file, err := parse.Parse(fileName, f)
	if err != nil {
		return err
	}

	err = file.ReplaceWithImporter()
	if err != nil {
		return err
	}

	return nil
}

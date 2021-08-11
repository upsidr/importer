package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/upsidr/importer/internal/parse"
)

var (
	updateCmd = &cobra.Command{
		Aliases: []string{"up"},
		Use:     "update [filename]",
		Short:   "Processes Importer markers and update the file in place",
		Long: `
` + "`update`" + ` command parses the provided file and processes the Import markers in place.

This does not support creating a new file, nor send the result to stdout. For such use cases, use ` + "`generate`" + ` command
`,
		RunE: executeUpdate,
	}
)

func executeUpdate(cmd *cobra.Command, args []string) error {
	// TODO: add some util func to hande all common error cases
	if len(args) != 1 {
		return errors.New("error: incorrect argument, you can only pass in 1 argument")
	}

	arg := args[0]
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

	err = file.ProcessMarkers()
	if err != nil {
		return err
	}

	err = file.ReplaceWithAfter()
	if err != nil {
		return err
	}

	return nil
}

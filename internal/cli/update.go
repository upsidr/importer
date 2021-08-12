package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/upsidr/importer/internal/file"
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

func init() {
	updateCmd.Flags().BoolVar(&isDryRun, "dry-run", false, "Run without updating the file")
}

func executeUpdate(cmd *cobra.Command, args []string) error {
	// TODO: add some util func to hande all common error cases

	if len(args) < 1 {
		return errors.New("missing file input")
	}

	for _, file := range args {
		if err := update(file); err != nil {
			fmt.Printf("Warning: failed to generate for '%s', %v", file, err)
		}
	}

	return nil
}

func update(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := parse.Parse(fileName, f)
	if err != nil {
		return err
	}

	err = fi.ProcessMarkers()
	if err != nil {
		return err
	}

	switch {
	case isDryRun:
		err = fi.ReplaceWithAfter(file.WithDryRun())
	default:
		err = fi.ReplaceWithAfter()
	}
	if err != nil {
		return err
	}

	return nil
}

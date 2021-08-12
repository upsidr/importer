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
	purgeCliCmd = &cobra.Command{
		Use:   "purge [filename]",
		Short: "Removes all imported lines and update the file in place",
		Long: `
` + "`purge`" + ` command processes the provided file and removes all the contents surrounded by Importer markers.

Importer markers will be left intact.
`,
		RunE: executePurge,
	}
)

func init() {
	purgeCliCmd.Flags().BoolVar(&isDryRun, "dry-run", false, "Run without updating the file")
}

func executePurge(cmd *cobra.Command, args []string) error {
	// TODO: add some util func to hande all common error cases

	if len(args) < 1 {
		return errors.New("missing file input")
	}

	for _, file := range args {
		if err := purge(file); err != nil {
			fmt.Printf("Warning: failed to purge for '%s', %v", file, err)
		}
	}

	return nil
}

func purge(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := parse.Parse(fileName, f)
	if err != nil {
		return err
	}

	switch {
	case isDryRun:
		err = fi.ReplaceWithPurged(file.WithDryRun())
	default:
		err = fi.ReplaceWithPurged()
	}
	if err != nil {
		return err
	}

	return nil
}

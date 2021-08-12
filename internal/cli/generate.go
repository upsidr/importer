package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/upsidr/importer/internal/parse"
)

var (
	generateCliCmd = &cobra.Command{
		Aliases: []string{"gen"},
		Use:     "generate [filename]",
		Short:   "Processes Importer markers and send output to stdout or file",
		Long: `
` + "`generate`" + ` command parses the provided file as the input, and output the processed file content to stdout or a file.

While ` + "`update`" + ` command is useful for managing file content in itself, ` + "`generate`" + ` can be used to create a separate template file.
This approach allows the input file to be full of Importer markes without actual importing, and only used as the template to generate a new file.
`,
		Args: cobra.MinimumNArgs(1),
		RunE: executeGenerate,
	}
	generateTargetFile string
)

func init() {
	generateCliCmd.Flags().StringVarP(&generateTargetFile, "out", "o", "", "write to `FILE`")
}

func executeGenerate(cmd *cobra.Command, args []string) error {
	// TODO: add some util func to hande all common error cases

	if len(args) < 1 {
		return errors.New("incorrect argument, you need to pass in an argument")
	}

	// Suppress usage message after this point
	cmd.SilenceUsage = true

	arg := args[0]
	out := generateTargetFile
	if err := generate(arg, out); err != nil {
		return fmt.Errorf("handling generate, %v", err)
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

	err = file.ProcessMarkers()
	if err != nil {
		return err
	}

	if targetFilepath != "" {
		return file.WriteAfterTo(targetFilepath)
	}

	return file.PrintAfter()
}

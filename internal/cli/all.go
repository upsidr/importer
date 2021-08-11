package cli

import (
	"github.com/spf13/cobra"
)

var (
	rootCmdName = "importer"
)

func Run(args []string) error {
	cmd := &cobra.Command{
		Use:   rootCmdName + " [command]",
		Short: "Import any lines, from anywhere",
		// Long: "To be updated",
	}
	cmd.AddCommand(
		previewCliCmd,
		updateCmd,
		generateCliCmd,
		purgeCliCmd,
	)
	return cmd.Execute()
}

package cli

import (
	"github.com/spf13/cobra"
	"github.com/urfave/cli/v2"
)

var (
	rootCmdName = "importer"

	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "importer",
		Short: "Code generation for any file with importer annotation",
		Long: `Importer allows pulling in any content from any file.

This is especially useful for file format that requires single file input, such as YAML and Markdown.
Within those files, you can add importer annotation to pull some content from other file.
As long as you have some code generation / compilation logic built into the CI setup,
you don't have to duplicate content in Markdowns, YAMLs, or any other files.`,
	}
)

// Execute executes the root command.
func Execute(args []string) error {
	rootCmd.AddCommand(
		previewCmd,
		generateCmd,
		purgeCmd,
	)
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func Run(args []string) error {
	app := &cli.App{
		Usage:                "Import any lines, from anywhere",
		UsageText:            "importer [command]",
		EnableBashCompletion: true,
	}
	app.Commands = []*cli.Command{
		previewCliCmd,
		updateCmd,
		generateCliCmd,
		purgeCliCmd,
	}
	return app.Run(args)
}

package cli

import (
	"github.com/spf13/cobra"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string

	rootCmd = &cobra.Command{
		Use:   "importer",
		Short: "Code generation for any file with importer annotation",
		Long: `Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		Run: func(cmd *cobra.Command, args []string) {
			//
		},
	}
)

// Execute executes the root command.
func Execute(args []string) error {
	rootCmd.AddCommand(
		previewCmd,
		generateCmd,
	)
	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

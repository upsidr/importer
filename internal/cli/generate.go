package cli

import (
	"github.com/spf13/cobra"
)

var (
	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Parse the provided file and generate result with imported files",
		Run: func(cmd *cobra.Command, args []string) {
			//
		},
	}
)

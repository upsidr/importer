package cli

import (
	"github.com/spf13/cobra"
)

var (
	previewCmd = &cobra.Command{
		Use:   "preview",
		Short: "Parse the provided file and send the result to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			//
		},
	}
)

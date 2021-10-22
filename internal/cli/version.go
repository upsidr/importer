package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/upsidr/importer/internal/version"
)

var (
	versionCmd = &cobra.Command{
		Aliases: []string{"v"},
		Use:     "version",
		Short:   "Print version information",
		Run:     executeVersion,
	}
)

func executeVersion(cmd *cobra.Command, args []string) {
	// Suppress usage message after this point
	cmd.SilenceUsage = true

	v := version.GetVersion()
	fmt.Println(v.VersionInfo())
}

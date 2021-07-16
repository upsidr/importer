package cli

import (
	"github.com/urfave/cli/v2"
)

var (
	rootCmdName = "importer"
)

func Run(args []string) error {
	app := &cli.App{
		Usage:                "Import any lines, from anywhere",
		UsageText:            rootCmdName + " [command]",
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

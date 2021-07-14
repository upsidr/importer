package main

import (
	"fmt"
	"os"

	"github.com/upsidr/importer/internal/cli"
)

func main() {
	// err := cli.Execute(os.Args[1:])
	err := cli.Run(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

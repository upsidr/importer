package main

import (
	"os"

	"github.com/upsidr/importer/internal/cli"
)

func main() {
	err := cli.Run(os.Args)
	if err != nil {
		// fmt.Fprintln(os.Stderr, err) // With cobra, the error is printed out already
		os.Exit(1)
	}
}

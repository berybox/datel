package main

import (
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	cliApp := &cli.App{
		UseShortOptionHandling: false,
		Usage:                  "Datel - database manipulation tool",
		Commands: []*cli.Command{
			&commandRunServer,
			&commandAddAdmin,
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

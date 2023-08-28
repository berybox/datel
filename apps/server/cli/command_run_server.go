package main

import (
	"os"

	"github.com/berybox/datel/apps/server/server"
	"github.com/urfave/cli/v2"
)

var (
	flagURI          = &cli.StringFlag{Name: "uri", Aliases: []string{"u"}, Usage: "MongoDB connection URI", Required: true}
	flagAddress      = &cli.StringFlag{Name: "address", Aliases: []string{"a"}, Usage: "Address on which this server will run", Value: ":4158"}
	flagUserOverride = &cli.StringFlag{Name: "overrideuser", Aliases: []string{"o"}, Usage: "If set, the server will always log in with this user ID. This feature is helpful for local use"}

	commandRunServer = cli.Command{
		Name:                   "run",
		Usage:                  "Run Datel server",
		UseShortOptionHandling: false,
		Flags: []cli.Flag{
			flagURI,
			flagAddress,
			flagUserOverride,
		},
		Action: func(ctx *cli.Context) error {
			return server.Run(
				ctx.String(flagURI.Name),
				ctx.String(flagAddress.Name),
				ctx.String(flagUserOverride.Name),
				os.Getenv("GODEBUGMODE") == "1",
			)
		},
	}
)

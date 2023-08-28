package main

import "github.com/urfave/cli/v2"

var (
	flagUsername = &cli.StringFlag{Name: "username", Aliases: []string{"s"}, Usage: "Name of user", Required: true}
	flagUserID   = &cli.StringFlag{Name: "id", Aliases: []string{"i"}, Usage: "User ID - must be unique and contain only letters and numbers", Required: true}

	commandAddAdmin = cli.Command{
		Name:                   "add-admin",
		Usage:                  "Add admin account",
		UseShortOptionHandling: false,
		Flags: []cli.Flag{
			flagURI,
			flagUsername,
			flagUserID,
		},
		Action: func(ctx *cli.Context) error {
			return addAdmin(
				ctx.String(flagURI.Name),
				ctx.String(flagUsername.Name),
				ctx.String(flagUserID.Name),
			)
		},
	}
)

package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/command"
)

var GlobalFlags = []cli.Flag{}

var Commands = []cli.Command{
	{
		Name:  "create",
		Usage: "",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "item",
				Usage:  "Create a new item in Qiita.",
				Action: command.CreateItem,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "file, f",
						Usage: "Path to the markdown file to item",
					},
				},
			},
		},
	},
	{
		Name:  "show",
		Usage: "",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "item",
				Usage:  "",
				Action: command.ShowItem,
				Flags:  []cli.Flag{},
			},
			{
				Name:   "items",
				Usage:  "",
				Action: command.ShowItems,
				Flags:  []cli.Flag{},
			},
		},
	},
	{
		Name:  "fetch",
		Usage: "",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "item",
				Usage:  "",
				Action: command.FetchItem,
				Flags:  []cli.Flag{},
			},
			{
				Name:   "items",
				Usage:  "",
				Action: command.FetchItems,
				Flags:  []cli.Flag{},
			},
		},
	},
	{
		Name:  "update",
		Usage: "",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "item",
				Usage:  "",
				Action: command.UpdateItem,
				Flags:  []cli.Flag{},
			},
			{
				Name:   "items",
				Usage:  "",
				Action: command.UpdateItems,
				Flags:  []cli.Flag{},
			},
		},
	},
	{
		Name:  "delete",
		Usage: "",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "items",
				Usage:  "",
				Action: command.DeleteItem,
				Flags:  []cli.Flag{},
			},
		},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

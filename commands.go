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
				Name:   "post",
				Usage:  "Create a new post in Qiita.",
				Action: command.CreatePost,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "file, f",
						Usage: "Path to the markdown file to post",
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
				Name:   "post",
				Usage:  "",
				Action: command.ShowPost,
				Flags:  []cli.Flag{},
			},
			{
				Name:   "posts",
				Usage:  "",
				Action: command.ShowPosts,
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
				Name:   "post",
				Usage:  "",
				Action: command.FetchPost,
				Flags:  []cli.Flag{},
			},
			{
				Name:   "posts",
				Usage:  "",
				Action: command.FetchPosts,
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
				Name:   "post",
				Usage:  "",
				Action: command.UpdatePost,
				Flags:  []cli.Flag{},
			},
			{
				Name:   "posts",
				Usage:  "",
				Action: command.UpdatePosts,
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
				Name:   "posts",
				Usage:  "",
				Action: command.DeletePost,
				Flags:  []cli.Flag{},
			},
		},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

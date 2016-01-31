package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/command"
)

var GlobalFlags = []cli.Flag{
	cli.BoolFlag{
		Name:  "debug, d",
		Usage: "Panic when error occurs",
	},
}

var Commands = []cli.Command{
	{
		Name:  "generate",
		Usage: "Generate something in your local",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "file",
				Usage:  "Generate a new markdown file for a new post",
				Action: command.GenerateFile,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "title, t",
						Usage: "The title of a new post",
					},
					cli.StringFlag{
						Name:  "team, T",
						Usage: "The name of a team, when you post to the team",
					},
				},
			},
		},
	},
	{
		Name:  "show",
		Usage: "Display resources",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "post",
				Usage:  "Print detail of a post in Qitta",
				Action: command.ShowPost,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Usage: "The ID of the post to be printed detail",
					},
				},
			},
			{
				Name:   "posts",
				Usage:  "Print a list of posts in Qiita",
				Action: command.ShowPosts,
				Flags:  []cli.Flag{},
			},
		},
	},
	{
		Name:  "pull",
		Usage: "Download resources from Qiita to current working directory",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "post",
				Usage:  "Download a post as a file",
				Action: command.PullPost,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Usage: "The ID of the post to be downloaded",
					},
				},
			},
			{
				Name:   "posts",
				Usage:  "Download posts as files",
				Action: command.PullPosts,
				Flags:  []cli.Flag{},
			},
		},
	},
	{
		Name:  "push",
		Usage: "Upload resources from current working directory to Qiita",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "post",
				Usage:  "Upload a post",
				Action: command.PushPost,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "filename, f",
						Usage: "The filename of the post to be uploaded",
					},
				},
			},
			// {
			// 	Name:   "posts",
			// 	Usage:  "Upload posts",
			// 	Action: command.PushPosts,
			// 	Flags:  []cli.Flag{},
			// },
		},
	},
	{
		Name:  "delete",
		Usage: "Delete resources from Qiita",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "post",
				Usage:  "Delete a post",
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

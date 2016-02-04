package cli

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/command"
	"github.com/minodisk/qiitactl/info"
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
				Action: command.CmdGenerateFile,
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
				Action: command.CmdShowPost,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Usage: "The ID of the post to be printed detail",
					},
					cli.StringFlag{
						Name:  "filename, f",
						Usage: "The filename of the post to be created",
					},
				},
			},
			{
				Name:   "posts",
				Usage:  "Print a list of posts in Qiita",
				Action: command.CmdShowPosts,
				Flags:  []cli.Flag{},
			},
		},
	},
	{
		Name:  "fetch",
		Usage: "Download resources from Qiita to current working directory",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "post",
				Usage:  "Download a post as a file",
				Action: command.CmdFetchPost,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "id, i",
						Usage: "The ID of the post to be downloaded",
					},
					cli.StringFlag{
						Name:  "filename, f",
						Usage: "The filename of the post to be created",
					},
				},
			},
			{
				Name:   "posts",
				Usage:  "Download posts as files",
				Action: command.CmdFetchPosts,
				Flags:  []cli.Flag{},
			},
		},
	},
	{
		Name:  "create",
		Usage: "Create resources from current working directory to Qiita",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "post",
				Usage:  "Create a post",
				Action: command.CmdCreatePost,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "filename, f",
						Usage: "The filename of the post to be created",
					},
					cli.BoolFlag{
						Name:  "tweet, t",
						Usage: "Tweet the post",
					},
					cli.BoolFlag{
						Name:  "gist, g",
						Usage: "Create codes in the post to GitHub Gist",
					},
				},
			},
		},
	},
	{
		Name:  "update",
		Usage: "Update resources from current working directory to Qiita",
		Flags: []cli.Flag{},
		Subcommands: []cli.Command{
			{
				Name:   "post",
				Usage:  "Update a post",
				Action: command.CmdUpdatePost,
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "filename, f",
						Usage: "The filename of the post to be updated",
					},
				},
			},
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

func NewApp() (app *cli.App) {
	app = cli.NewApp()
	app.Name = info.Name
	app.Version = info.Version
	app.Author = info.Author
	app.Usage = "Controls the Qiita posts"
	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound
	return
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}

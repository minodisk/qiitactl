package command

import (
	"fmt"
	"io"

	"github.com/alecthomas/kingpin"
	"github.com/minodisk/qiitactl/api"
)

type Command struct {
	Client api.Client
	Out    io.Writer
	Error  io.Writer

	Application  *kingpin.Application
	Generate     *kingpin.CmdClause
	GenerateFile *kingpin.CmdClause
	Create       *kingpin.CmdClause
	CreatePost   *kingpin.CmdClause
	Show         *kingpin.CmdClause
	ShowPost     *kingpin.CmdClause
	ShowPosts    *kingpin.CmdClause
	Fetch        *kingpin.CmdClause
	FetchPost    *kingpin.CmdClause
	FetchPosts   *kingpin.CmdClause
	Update       *kingpin.CmdClause
	UpdatePost   *kingpin.CmdClause
	Delete       *kingpin.CmdClause
	DeletePost   *kingpin.CmdClause

	GlobalOptions      GlobalOptions
	GenerateFileRunner GenerateFileRunner
	CreatePostRunner   CreatePostRunner
	ShowPostRunner     ShowPostRunner
	ShowPostsRunner    ShowPostsRunner
	FetchPostRunner    FetchPostRunner
	FetchPostsRunner   FetchPostsRunner
	UpdatePostRunner   UpdatePostRunner
	DeletePostRunner   DeletePostRunner
}

type GlobalOptions struct {
	Debug *bool
}

func New(client api.Client, out io.Writer, err io.Writer) (c Command) {
	c.Client = client
	c.Out = out
	c.Error = err

	c.Application = kingpin.New("qiitactl", "A command-line chat application.")
	c.GlobalOptions = GlobalOptions{
		Debug: c.Application.Flag("debug", "Enable debug mode.").Bool(),
	}

	c.Generate = c.Application.Command("generate", "Generate something in your local.")
	c.GenerateFile = c.Generate.Command("file", "Generate a new markdown file for a new post.")
	c.GenerateFileRunner = GenerateFileRunner{
		Title: c.GenerateFile.Arg("title", "The title of a new post.").Required().String(),
		Team:  c.GenerateFile.Flag("team", "The name of a team, when you post to the team.").Short('t').String(),
	}

	c.Create = c.Application.Command("create", "Create resources from current working directory to Qiita.")
	c.CreatePost = c.Create.Command("post", "Create a post in Qiita.")
	c.CreatePostRunner = CreatePostRunner{
		File:  c.CreatePost.Arg("filename", "The filename of the post to be created").Required().File(),
		Tweet: c.CreatePost.Flag("tweet", "Tweet the created post in Twitter.").Short('t').Bool(),
		Gist:  c.CreatePost.Flag("gist", "Post codes in the created post to GitHub Gist.").Short('g').Bool(),
	}

	c.Show = c.Application.Command("show", "Display resources.")
	c.ShowPost = c.Show.Command("post", "Display detail of a post in Qitta.")
	c.ShowPostRunner = ShowPostRunner{
		ID:   c.ShowPost.Flag("id", "The ID of the post to be printed detail.").Short('i').String(),
		File: c.ShowPost.Flag("filename", "The filename of the post to be created.").Short('f').File(),
	}
	c.ShowPosts = c.Show.Command("posts", "Display posts in Qiita.")
	c.ShowPostsRunner = ShowPostsRunner{}

	c.Fetch = c.Application.Command("fetch", "Download resources from Qiita to current working directory.")
	c.FetchPost = c.Fetch.Command("post", "Download a post as a file.")
	c.FetchPostRunner = FetchPostRunner{
		ID:   c.FetchPost.Flag("id", "The ID of the post to be downloaded.").Short('i').String(),
		File: c.FetchPost.Flag("filename", "The filename of the post to be created.").Short('f').File(),
	}
	c.FetchPosts = c.Fetch.Command("posts", "Download posts as files.")
	c.FetchPostsRunner = FetchPostsRunner{}

	c.Update = c.Application.Command("update", "Update resources from current working directory to Qiita.")
	c.UpdatePost = c.Update.Command("post", "Update a post in Qiita.")
	c.UpdatePostRunner = UpdatePostRunner{
		File: c.UpdatePost.Arg("filename", "The filename of the post to be updated.").Required().File(),
	}

	c.Delete = c.Application.Command("delete", "Delete resources from current working directory to Qiita.")
	c.DeletePost = c.Delete.Command("post", "Delete a post in Qiita.")
	c.DeletePostRunner = DeletePostRunner{
		File: c.DeletePost.Arg("filename", "The filename of the post to be deleted.").Required().File(),
	}

	return
}

func (c Command) Run(args []string) {
	var err error

	switch kingpin.MustParse(c.Application.Parse(args[1:])) {
	case c.GenerateFile.FullCommand():
		err = c.GenerateFileRunner.Run(c.Client, c.GlobalOptions, c.Out)
	case c.CreatePost.FullCommand():
		err = c.CreatePostRunner.Run(c.Client, c.GlobalOptions, c.Out)
	case c.ShowPost.FullCommand():
		err = c.ShowPostRunner.Run(c.Client, c.GlobalOptions, c.Out)
	case c.ShowPosts.FullCommand():
		err = c.ShowPostsRunner.Run(c.Client, c.GlobalOptions, c.Out)
	case c.FetchPost.FullCommand():
		err = c.FetchPostRunner.Run(c.Client, c.GlobalOptions, c.Out)
	case c.FetchPosts.FullCommand():
		err = c.FetchPostsRunner.Run(c.Client, c.GlobalOptions, c.Out)
	case c.UpdatePost.FullCommand():
		err = c.UpdatePostRunner.Run(c.Client, c.GlobalOptions, c.Out)
	case c.DeletePost.FullCommand():
		err = c.DeletePostRunner.Run(c.Client, c.GlobalOptions, c.Out)
	}

	if err != nil {
		// log.Println("--------------------")
		fmt.Fprintf(c.Error, "%s\n", err)
	}
}

// // GenerateApp generates an application of command line interface.
// func GenerateApp(client api.Client, outWriter io.Writer, errWriter io.Writer) (c *cli.App) {
// 	c = cli.NewApp()
// 	c.Writer = outWriter
// 	c.Name = info.Name
// 	c.Version = info.Version
// 	c.Author = info.Author
// 	c.Usage = "Controls the Qiita posts"
// 	c.Flags = []cli.Flag{
// 		cli.BoolFlag{
// 			Name:  "debug, d",
// 			Usage: "Panic when error occurs",
// 		},
// 	}
// 	c.Commands = []cli.Command{
// 		{
// 			Name:  "generate",
// 			Usage: "Generate something in your local",
// 			Flags: []cli.Flag{},
// 			Subcommands: []cli.Command{
// 				{
// 					Name:   "file",
// 					Usage:  "Generate a new markdown file for a new post",
// 					Action: partializeWithWriter(command.GenerateFile, client, outWriter, errWriter),
// 					Flags: []cli.Flag{
// 						cli.StringFlag{
// 							Name:  "title, t",
// 							Usage: "The title of a new post",
// 						},
// 						cli.StringFlag{
// 							Name:  "team, T",
// 							Usage: "The name of a team, when you post to the team",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Name:  "create",
// 			Usage: "Create resources from current working directory to Qiita",
// 			Flags: []cli.Flag{},
// 			Subcommands: []cli.Command{
// 				{
// 					Name:   "post",
// 					Usage:  "Create a post",
// 					Action: partialize(command.CreatePost, client, errWriter),
// 					Flags: []cli.Flag{
// 						cli.StringFlag{
// 							Name:  "filename, f",
// 							Usage: "The filename of the post to be created",
// 						},
// 						cli.BoolFlag{
// 							Name:  "tweet, t",
// 							Usage: "Tweet the created post in Twitter",
// 						},
// 						cli.BoolFlag{
// 							Name:  "gist, g",
// 							Usage: "Post codes in the created post to GitHub Gist",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Name:  "show",
// 			Usage: "Display resources",
// 			Flags: []cli.Flag{},
// 			Subcommands: []cli.Command{
// 				{
// 					Name:   "post",
// 					Usage:  "Print detail of a post in Qitta",
// 					Action: partializeWithWriter(command.ShowPost, client, outWriter, errWriter),
// 					Flags: []cli.Flag{
// 						cli.StringFlag{
// 							Name:  "id, i",
// 							Usage: "The ID of the post to be printed detail",
// 						},
// 						cli.StringFlag{
// 							Name:  "filename, f",
// 							Usage: "The filename of the post to be created",
// 						},
// 					},
// 				},
// 				{
// 					Name:   "posts",
// 					Usage:  "Print a list of posts in Qiita",
// 					Action: partializeWithWriter(command.ShowPosts, client, outWriter, errWriter),
// 					Flags:  []cli.Flag{},
// 				},
// 			},
// 		},
// 		{
// 			Name:  "fetch",
// 			Usage: "Download resources from Qiita to current working directory",
// 			Flags: []cli.Flag{},
// 			Subcommands: []cli.Command{
// 				{
// 					Name:   "post",
// 					Usage:  "Download a post as a file",
// 					Action: partialize(command.FetchPost, client, errWriter),
// 					Flags: []cli.Flag{
// 						cli.StringFlag{
// 							Name:  "id, i",
// 							Usage: "The ID of the post to be downloaded",
// 						},
// 						cli.StringFlag{
// 							Name:  "filename, f",
// 							Usage: "The filename of the post to be created",
// 						},
// 					},
// 				},
// 				{
// 					Name:   "posts",
// 					Usage:  "Download posts as files",
// 					Action: partialize(command.FetchPosts, client, errWriter),
// 					Flags:  []cli.Flag{},
// 				},
// 			},
// 		},
// 		{
// 			Name:  "update",
// 			Usage: "Update resources from current working directory to Qiita",
// 			Flags: []cli.Flag{},
// 			Subcommands: []cli.Command{
// 				{
// 					Name:   "post",
// 					Usage:  "Update a post",
// 					Action: partialize(command.UpdatePost, client, errWriter),
// 					Flags: []cli.Flag{
// 						cli.StringFlag{
// 							Name:  "filename, f",
// 							Usage: "The filename of the post to be updated",
// 						},
// 					},
// 				},
// 			},
// 		},
// 		{
// 			Name:  "delete",
// 			Usage: "Delete resources from Qiita",
// 			Flags: []cli.Flag{},
// 			Subcommands: []cli.Command{
// 				{
// 					Name:   "post",
// 					Usage:  "Delete a post",
// 					Action: partialize(command.DeletePost, client, errWriter),
// 					Flags: []cli.Flag{
// 						cli.StringFlag{
// 							Name:  "filename, f",
// 							Usage: "The filename of the post to be updated",
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// 	c.CommandNotFound = func(c *cli.Context, command string) {
// 		fmt.Fprintf(errWriter, "%s: '%s' is not a %s command. See '%s --help'.\n", c.App.Name, command, c.App.Name, c.App.Name)
// 		// os.Exit(2)
// 	}
// 	return
// }
//
// func partialize(cmdFunc func(*cli.Context, api.Client) error, client api.Client, errWriter io.Writer) func(ctx *cli.Context) {
// 	return func(ctx *cli.Context) {
// 		client.DebugMode(ctx.GlobalBool("debug"))
// 		err := cmdFunc(ctx, client)
// 		if err != nil {
// 			printError(ctx, errWriter, err)
// 			// os.Exit(2)
// 		}
// 	}
// }
//
// func partializeWithWriter(cmdFunc func(*cli.Context, api.Client, io.Writer) error, client api.Client, outWriter io.Writer, errWriter io.Writer) func(ctx *cli.Context) {
// 	return func(ctx *cli.Context) {
// 		client.DebugMode(ctx.GlobalBool("debug"))
// 		err := cmdFunc(ctx, client, outWriter)
// 		if err != nil {
// 			printError(ctx, errWriter, err)
// 			// os.Exit(2)
// 		}
// 	}
// }
//
// func printError(ctx *cli.Context, w io.Writer, err error) {
// 	fmt.Fprint(w, err)
// }

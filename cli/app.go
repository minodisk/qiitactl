package cli

import (
	"io"

	"github.com/alecthomas/kingpin"
	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/api"
)

var (
	app   = kingpin.New("qiitactl", "A command-line chat application.")
	debug = app.Flag("debug", "Enable debug mode.").Bool()

	generate          = app.Command("generate", "Generate something in your local.")
	generateFile      = generate.Command("file", "Generate a new markdown file for a new post.")
	generateFileTitle = generateFile.Arg("title", "The title of a new post.").Required().String()
	generateFileTeam  = generateFile.Flag("team", "The name of a team, when you post to the team.").Short('t').String()

	create          = app.Command("create", "Create resources from current working directory to Qiita.")
	createPost      = create.Command("post", "Create a post in Qiita.")
	createPostFile  = createPost.Arg("filename", "The filename of the post to be created").Required().File()
	createPostTweet = createPost.Flag("tweet", "Tweet the created post in Twitter.").Short('t').Bool()
	createPostGist  = createPost.Flag("gist", "Post codes in the created post to GitHub Gist.").Short('g').Bool()

	show         = app.Command("show", "Display resources.")
	showPost     = show.Command("post", "Display detail of a post in Qitta.")
	showPostID   = showPost.Flag("id", "The ID of the post to be printed detail.").Short('i').String()
	showPostFile = showPost.Flag("filename", "The filename of the post to be created.").Short('f').File()
	showPosts    = show.Command("posts", "Display posts in Qiita.")

	fetch         = app.Command("fetch", "Download resources from Qiita to current working directory.")
	fetchPost     = fetch.Command("post", "Download a post as a file.")
	fetchPostID   = fetchPost.Flag("id", "The ID of the post to be downloaded.").Short('i').String()
	fetchPostFile = fetchPost.Flag("filename", "The filename of the post to be created.").Short('f').File()
	fetchPosts    = fetch.Command("posts", "Download posts as files.")

	update         = app.Command("update", "Update resources from current working directory to Qiita.")
	updatePost     = update.Command("post", "Update a post in Qiita.")
	updatePostFile = updatePost.Arg("filename", "The filename of the post to be updated.").Required().File()

	deleteCmd      = app.Command("delete", "Delete resources from current working directory to Qiita.")
	deletePost     = deleteCmd.Command("delete", "Delete a post in Qiita.")
	deletePostFile = deletePost.Arg("filename", "The filename of the post to be deleted.").Required().File()
)

func GenerateApp(client api.Client, outWriter io.Writer, errWriter io.Writer) (app *cli.App) {
	// switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	switch kingpin.Parse() {
	// Register user
	case show.FullCommand():
		// println(*registerNick)
		// Post message
		// case post.FullCommand():
		// 	if *postImage != nil {
		// 	}
		// 	text := strings.Join(*postText, " ")
		// 	println("Post:", text)
	}
}

// // GenerateApp generates an application of command line interface.
// func GenerateApp(client api.Client, outWriter io.Writer, errWriter io.Writer) (app *cli.App) {
// 	app = cli.NewApp()
// 	app.Writer = outWriter
// 	app.Name = info.Name
// 	app.Version = info.Version
// 	app.Author = info.Author
// 	app.Usage = "Controls the Qiita posts"
// 	app.Flags = []cli.Flag{
// 		cli.BoolFlag{
// 			Name:  "debug, d",
// 			Usage: "Panic when error occurs",
// 		},
// 	}
// 	app.Commands = []cli.Command{
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
// 	app.CommandNotFound = func(c *cli.Context, command string) {
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

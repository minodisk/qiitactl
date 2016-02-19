package command

import (
	"fmt"
	"io"

	"github.com/alecthomas/kingpin"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/info"
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

func New(info info.Info, client api.Client, out io.Writer, err io.Writer) (c Command) {
	c.Client = client
	c.Out = out
	c.Error = err

	c.Application = kingpin.New("qiitactl", "Command line interface to manage the posts in Qitta.")
	c.Application.Version(info.Version)
	c.Application.Author(info.Author)
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

	cmd, err := c.Application.Parse(args[1:])
	c.Client.DebugMode(*c.GlobalOptions.Debug)

	switch kingpin.MustParse(cmd, err) {
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
		fmt.Fprintf(c.Error, "%s\n", err)
	}
}

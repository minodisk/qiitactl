package command

import (
	"fmt"
	"io"
	"os"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

type ShowPostRunner struct {
	ID   *string
	File **os.File
}

// ShowPost outputs your post fetched from Qiita to stdout.
func (r ShowPostRunner) Run(c api.Client, o GlobalOptions, w io.Writer) (err error) {
	id, err := getID(*r.ID, *r.File)
	if err != nil {
		return
	}
	post, err := model.FetchPost(c, nil, id)
	if err != nil {
		return
	}
	err = printPost(w, post)
	return
}

type ShowPostsRunner struct{}

// ShowPosts outputs your posts fetched from Qiita to stdout.
func (r ShowPostsRunner) Run(c api.Client, o GlobalOptions, w io.Writer) (err error) {
	posts, err := model.FetchPosts(c, nil)
	if err != nil {
		return
	}
	err = printPosts(w, posts, nil)
	if err != nil {
		return
	}

	teams, err := model.FetchTeams(c)
	if err != nil {
		return
	}
	for _, team := range teams {
		posts, err = model.FetchPosts(c, &team)
		if err != nil {
			return
		}
		err = printPosts(w, posts, &team)
		if err != nil {
			return
		}
	}
	return
}

func printPosts(w io.Writer, posts model.Posts, team *model.Team) (err error) {
	if team == nil {
		_, err = w.Write([]byte("Posts in Qiita:\n"))
	} else {
		_, err = w.Write([]byte(fmt.Sprintf("Posts in Qiita:Team (%s):\n", team.Name)))
	}
	if err != nil {
		return
	}
	for _, post := range posts {
		err = printPost(w, post)
		if err != nil {
			return
		}
	}
	return
}

func printPost(w io.Writer, post model.Post) (err error) {
	_, err = w.Write([]byte(fmt.Sprintf("%s %s %s\n", post.ID, post.CreatedAt.FormatDate(), post.Title)))
	return
}

type FetchPostRunner struct {
	ID   *string
	File **os.File
}

// FetchPost fetches your post from Qiita to current working directory.
func (r FetchPostRunner) Run(c api.Client, o GlobalOptions, w io.Writer) (err error) {
	id, err := getID(*r.ID, *r.File)
	if err != nil {
		return
	}
	post, err := model.FetchPost(c, nil, id)
	if err != nil {
		return
	}
	err = post.Save(nil)
	return
}

func getID(id string, file *os.File) (i string, err error) {
	if id != "" {
		i = id
		return
	}
	if file != nil {
		post, err := model.NewPostWithOSFile(file)
		if err != nil {
			return "", err
		}
		i = post.ID
		return i, nil
	}

	err = fmt.Errorf("fetch post: id or filename is required")
	return
}

type FetchPostsRunner struct{}

// FetchPosts fetches your posts from Qiita to current working directory.
func (r FetchPostsRunner) Run(c api.Client, o GlobalOptions, w io.Writer) (err error) {
	posts, err := model.FetchPosts(c, nil)
	if err != nil {
		return
	}
	err = posts.Save()
	if err != nil {
		return
	}

	teams, err := model.FetchTeams(c)
	if err != nil {
		return
	}
	for _, team := range teams {
		var posts model.Posts
		posts, err = model.FetchPosts(c, &team)
		if err != nil {
			return
		}
		err = posts.Save()
		if err != nil {
			return
		}
	}
	return
}

type CreatePostRunner struct {
	File  **os.File
	Tweet *bool
	Gist  *bool
}

// CreatePost creates a new post in Qiita with a specified file.
func (r CreatePostRunner) Run(c api.Client, o GlobalOptions, w io.Writer) (err error) {
	opts := model.CreationOptions{
		Tweet: *r.Tweet,
		Gist:  *r.Gist,
	}

	post, err := model.NewPostWithOSFile(*r.File)
	if err != nil {
		return
	}
	err = post.Create(c, opts)
	if err != nil {
		return
	}
	err = post.Save(nil)
	return
}

type UpdatePostRunner struct {
	File **os.File
}

// UpdatePost updates your post in Qiita with a specified file.
func (r UpdatePostRunner) Run(c api.Client, o GlobalOptions, w io.Writer) (err error) {
	post, err := model.NewPostWithOSFile(*r.File)
	if err != nil {
		return
	}
	err = post.Update(c)
	if err != nil {
		return
	}
	err = post.Save(nil)
	return
}

type DeletePostRunner struct {
	File **os.File
}

// DeletePost deletes your post from Qiita with a specified file.
func (r DeletePostRunner) Run(c api.Client, o GlobalOptions, w io.Writer) (err error) {
	post, err := model.NewPostWithOSFile(*r.File)
	if err != nil {
		return
	}
	err = post.Delete(c)
	if err != nil {
		return
	}
	err = post.Save(nil)
	return
}

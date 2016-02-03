package command

import (
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

func CmdShowPost(c *cli.Context) {
	client, err := api.NewClient(nil)
	if err != nil {
		printError(c, err)
		return
	}

	err = ShowPost(client, os.Stdout, c.String("id"), c.String("filename"))
	if err != nil {
		printError(c, err)
	}
}

func ShowPost(client api.Client, w io.Writer, id, filename string) (err error) {
	id, err = getID(id, filename)
	if err != nil {
		return
	}
	post, err := model.FetchPost(client, nil, id)
	if err != nil {
		return
	}
	err = printPost(w, post)
	return
}

func CmdShowPosts(c *cli.Context) {
	client, err := api.NewClient(nil)
	if err != nil {
		printError(c, err)
		return
	}
	err = ShowPosts(client, os.Stdout)
	if err != nil {
		printError(c, err)
	}
}

func ShowPosts(client api.Client, w io.Writer) (err error) {
	posts, err := model.FetchPosts(client, nil)
	if err != nil {
		return
	}
	err = printPosts(w, posts, nil)
	if err != nil {
		return
	}

	teams, err := model.FetchTeams(client)
	if err != nil {
		return
	}
	for _, team := range teams {
		posts, err = model.FetchPosts(client, &team)
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

func CmdFetchPost(c *cli.Context) {
	client, err := api.NewClient(nil)
	if err != nil {
		printError(c, err)
		return
	}
	err = FetchPost(client, c.String("id"), c.String("filename"))
	if err != nil {
		printError(c, err)
	}
}

func FetchPost(client api.Client, id, filename string) (err error) {
	id, err = getID(id, filename)
	if err != nil {
		return
	}
	post, err := model.FetchPost(client, nil, id)
	if err != nil {
		return
	}
	err = post.Save()
	return
}
func CmdFetchPosts(c *cli.Context) {
	client, err := api.NewClient(nil)
	if err != nil {
		printError(c, err)
		return
	}
	err = FetchPosts(client)
	if err != nil {
		printError(c, err)
	}
}

func getID(id, filename string) (i string, err error) {
	if id != "" {
		i = id
		return
	}
	if filename == "" {
		err = fmt.Errorf("fetch post: id or filename is required")
		return
	}
	post, err := model.NewPostWithFile(filename)
	if err != nil {
		return
	}
	i = post.ID
	return
}

func FetchPosts(client api.Client) (err error) {
	posts, err := model.FetchPosts(client, nil)
	if err != nil {
		return
	}
	err = posts.Save()
	if err != nil {
		return
	}

	teams, err := model.FetchTeams(client)
	if err != nil {
		return
	}
	for _, team := range teams {
		var posts model.Posts
		posts, err = model.FetchPosts(client, &team)
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

func CmdCreatePost(c *cli.Context) {
	err := CreatePost(c.String("filename"))
	if err != nil {
		printError(c, err)
	}
}

func CreatePost(filename string) (err error) {
	post, err := model.NewPostWithFile(filename)
	if err != nil {
		return
	}
	err = post.Create()
	return
}

func CmdUpdatePost(c *cli.Context) {
	client, err := api.NewClient(nil)
	if err != nil {
		printError(c, err)
		return
	}
	err = UpdatePost(client, c.String("filename"))
	if err != nil {
		printError(c, err)
	}
}

func UpdatePost(client api.Client, filename string) (err error) {
	post, err := model.NewPostWithFile(filename)
	if err != nil {
		return
	}
	err = post.Update(client)
	return
}

func DeletePost(c *cli.Context) {
}

// func PostsDiff(commit1, commit2 string) (err error) {
// 	fmt.Printf("Post diff between %s and %s\n", commit1, commit2)
//
// 	err = exec.Command("git", "config", "--local", "core.quotepath", "false").Run()
// 	if err != nil {
// 		return
// 	}
//
// 	cmd := exec.Command("git", "--no-pager", "diff", "--name-only", commit1, commit2)
//
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return
// 	}
//
// 	filenames, err := cmd.Output()
// 	for _, filename := range strings.Split(string(filenames), "\n") {
// 		if filename == "" {
// 			continue
// 		}
// 		filename = filepath.Join(wd, strings.Trim(filename, "\""))
// 		fmt.Println(filename)
//
// 		b, err := ioutil.ReadFile(filename)
// 		if err != nil {
// 			return err
// 		}
// 		_, err = model.NewPost(string(b))
// 		if err != nil {
// 			return err
// 		}
// 	}
//
// 	return
// }

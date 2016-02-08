package command

import (
	"fmt"
	"io"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

func ShowPost(c *cli.Context, client api.Client, w io.Writer) (err error) {
	id := c.String("id")
	filename := c.String("filename")
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

func ShowPosts(c *cli.Context, client api.Client, w io.Writer) (err error) {
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

func FetchPost(c *cli.Context, client api.Client) (err error) {
	id := c.String("id")
	filename := c.String("filename")

	id, err = getID(id, filename)
	if err != nil {
		return
	}
	post, err := model.FetchPost(client, nil, id)
	if err != nil {
		return
	}
	err = post.Save(nil)
	return
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

func FetchPosts(c *cli.Context, client api.Client) (err error) {
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

func CreatePost(c *cli.Context, client api.Client) (err error) {
	filename := c.String("filename")
	opts := model.CreationOptions{
		Tweet: c.Bool("tweet"),
		Gist:  c.Bool("gist"),
	}

	post, err := model.NewPostWithFile(filename)
	if err != nil {
		return
	}
	err = post.Create(client, opts)
	if err != nil {
		return
	}
	err = post.Save(nil)
	return
}

func UpdatePost(c *cli.Context, client api.Client) (err error) {
	filename := c.String("filename")

	post, err := model.NewPostWithFile(filename)
	if err != nil {
		return
	}
	err = post.Update(client)
	if err != nil {
		return
	}
	err = post.Save(nil)
	return
}

func DeletePost(c *cli.Context, client api.Client) (err error) {
	filename := c.String("filename")

	post, err := model.NewPostWithFile(filename)
	if err != nil {
		return
	}
	err = post.Delete(client)
	if err != nil {
		return
	}
	err = post.Save(nil)
	return
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

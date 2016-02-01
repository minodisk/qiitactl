package command

import (
	"fmt"
	"io"
	"os"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

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
	posts, err := model.FetchPosts(client)
	if err != nil {
		return
	}
	printPosts(w, posts, nil)

	teams, err := model.FetchTeams(client)
	if err != nil {
		return
	}
	for _, team := range teams {
		posts, err = model.FetchPostsInTeam(client, &team)
		if err != nil {
			return
		}
		printPosts(w, posts, &team)
	}
	return
}

func printPosts(w io.Writer, posts model.Posts, team *model.Team) {
	if team == nil {
		w.Write([]byte("Posts in Qiita:\n"))
	} else {
		w.Write([]byte(fmt.Sprintf("Posts in Qiita:Team (%s):\n", team.Name)))
	}
	for _, post := range posts {
		w.Write([]byte(fmt.Sprintf("%s %s %s\n", post.Id, post.CreatedAt.FormatDate(), post.Title)))
	}
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

func FetchPosts(client api.Client) (err error) {
	posts, err := model.FetchPosts(client)
	if err != nil {
		return
	}
	err = posts.SaveToLocal()
	if err != nil {
		return
	}

	teams, err := model.FetchTeams(client)
	if err != nil {
		return
	}
	for _, team := range teams {
		var posts model.Posts
		posts, err = model.FetchPostsInTeam(client, &team)
		if err != nil {
			return
		}
		err = posts.SaveToLocal()
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
	file, err := model.NewFileFromLocal(filename)
	if err != nil {
		return
	}
	err = file.Post.Create()
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
	file, err := model.NewFileFromLocal(filename)
	if err != nil {
		return
	}
	err = file.Post.Update(client)
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

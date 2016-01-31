package command

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

func printError(c *cli.Context, err error) {
	if c.GlobalBool("debug") {
		panic(err)
	} else {
		fmt.Println(err)
	}
}

func ShowPosts(c *cli.Context) {
	err := func() (err error) {
		client, err := api.NewClient()
		if err != nil {
			return
		}

		posts, err := model.FetchPosts(client)
		if err != nil {
			return
		}
		printPosts(posts, nil)

		teams, err := model.FetchTeams(client)
		if err != nil {
			return
		}
		for _, team := range teams {
			posts, err = model.FetchPostsInTeam(client, &team)
			if err != nil {
				return
			}
			printPosts(posts, &team)
		}
		return
	}()
	if err != nil {
		printError(c, err)
	}
}

func printPosts(posts model.Posts, team *model.Team) {
	if team == nil {
		fmt.Println("Posts in Qiita:")
	} else {
		fmt.Printf("Posts in Qiita:Team (%s):\n", team.Name)
	}
	for _, post := range posts {
		fmt.Println(post.Id, post.CreatedAt.FormatDate(), post.Title)
	}
}

func ShowPost(c *cli.Context) {
}

func FetchPost(c *cli.Context) {
	// Write your code here
}

func FetchPosts(c *cli.Context) {
	err := func() (err error) {
		client, err := api.NewClient()
		if err != nil {
			return
		}

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
	}()
	if err != nil {
		printError(c, err)
	}
}

func CreatePost(c *cli.Context) {
	err := func() (err error) {
		filename := c.String("filename")
		file, err := model.NewFileFromLocal(filename)
		if err != nil {
			return
		}
		file.Post.Create()
		return
	}()
	if err != nil {
		printError(c, err)
	}
}

func CreatePosts(c *cli.Context) {
	// Write your code here
}

func DeletePost(c *cli.Context) {
	// Write your code here
}

// func fetchAllPosts() (err error) {
// 	err = fetchPosts("")
// 	if err != nil {
// 		return
// 	}
// 	err = fetchPosts("")
// 	return
// }

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

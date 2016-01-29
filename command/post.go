package command

import (
	"log"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/model"
)

func CreatePost(c *cli.Context) {
	err := func() (err error) {
		// client := model.NewClient()
		filename := c.String("file")
		_, err = model.NewPostFromFile(filename)
		if err != nil {
			return
		}
		return
	}()
	if err != nil {
		panic(err)
	}
}

func ShowPosts(c *cli.Context) {
	err := func() (err error) {
		client := model.NewClient()
		err = model.ShowPosts(client)
		return
	}()
	if err != nil {
		log.Fatal(err)
	}
}

func ShowPost(c *cli.Context) {
}

func FetchPost(c *cli.Context) {
	// Write your code here
}

func FetchPosts(c *cli.Context) {
	err := func() (err error) {
		client := model.NewClient()

		posts, err := model.FetchPosts(client)
		if err != nil {
			return
		}
		err = posts.SaveToLocal("mine")
		if err != nil {
			return
		}

		var teams model.Teams
		err = teams.Fetch(client)
		if err != nil {
			return
		}
		for _, team := range teams {
			var posts model.Posts
			posts, err = model.FetchPostsInTeam(client, team)
			if err != nil {
				return
			}
			err = posts.SaveToLocal(team.Name)
			if err != nil {
				return
			}
		}
		return
	}()
	if err != nil {
		panic(err)
	}
}

func UpdatePost(c *cli.Context) {
	// Write your code here
}

func UpdatePosts(c *cli.Context) {
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

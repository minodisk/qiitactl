package command

import (
	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/model"
)

func GenerateFile(c *cli.Context) {
	err := func() (err error) {
		teamName := c.String("team")
		var team *model.Team
		if teamName != "" {
			team.Name = teamName
		}

		post := model.NewPost()
		post.Title = c.String("title")
		post.Team = team
		file := model.NewFile(post)
		err = file.Save()
		return
	}()
	if err != nil {
		printError(c, err)
	}
}

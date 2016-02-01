package command

import (
	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/model"
)

func GenerateFile(c *cli.Context) {
	err := func() (err error) {
		teamID := c.String("team")
		var team *model.Team
		if teamID != "" {
			team.ID = teamID
		}

		post := model.NewPost()
		post.Title = c.String("title")
		post.Team = team
		err = post.Save()
		return
	}()
	if err != nil {
		printError(c, err)
	}
}

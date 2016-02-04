package command

import (
	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/model"
)

func GenerateFile(c *cli.Context) {
	err := func(teamID string, title string) (err error) {
		var team *model.Team
		if teamID != "" {
			team.ID = teamID
		}
		post := model.NewPost(title, nil, team)
		err = post.Save()
		return
	}(c.String("team"), c.String("title"))
	if err != nil {
		printError(c, err)
	}
}

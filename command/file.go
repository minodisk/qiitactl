package command

import (
	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

func GenerateFile(c *cli.Context, client api.Client) (err error) {
	teamID := c.String("team")
	title := c.String("title")

	var team *model.Team
	if teamID != "" {
		team.ID = teamID
	}
	post := model.NewPost(title, nil, team)
	err = post.Save()
	return
}

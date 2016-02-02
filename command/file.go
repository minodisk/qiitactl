package command

import (
	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/model"
)

func CmdGenerateFile(c *cli.Context) {
	err := GenerateFile(c.String("team"), c.String("title"))
	if err != nil {
		printError(c, err)
	}
}

func GenerateFile(teamID string, title string) (err error) {
	var team *model.Team
	if teamID != "" {
		team.ID = teamID
	}

	post := model.NewPost(title, team)

	err = post.Save()
	return
}

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
		file := model.NewFile(model.Post{
			Title: c.String("title"),
			Team:  team,
		})
		err = file.Save()
		return
	}()
	if err != nil {
		printError(c, err)
	}
}

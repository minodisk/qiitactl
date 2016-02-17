package command

import (
	"fmt"
	"io"

	"github.com/codegangsta/cli"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

// GenerateFile generates markdown file at current working directory.
func GenerateFile(c *cli.Context, client api.Client, w io.Writer) (err error) {
	teamID := c.String("team")
	title := c.String("title")

	var team *model.Team
	if teamID != "" {
		team = &model.Team{
			ID: teamID,
		}
	}
	post := model.NewPost(title, nil, team)
	err = post.Save(nil)
	if err != nil {
		return
	}
	_, err = fmt.Fprintf(w, "%s\n", post.Path)
	return
}

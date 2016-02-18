package command

import (
	"fmt"
	"io"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/model"
)

type GenerateFileRunner struct {
	Title *string
	Team  *string
}

// GenerateFile generates markdown file at current working directory.
func (r GenerateFileRunner) Run(c api.Client, o GlobalOptions, w io.Writer) (err error) {
	var team *model.Team
	if *r.Team != "" {
		team = &model.Team{
			ID: *r.Team,
		}
	}

	post := model.NewPost(*r.Title, nil, team)
	err = post.Save(nil)
	if err != nil {
		return
	}

	_, err = fmt.Fprintf(w, "%s\n", post.Path)

	return
}

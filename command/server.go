package command

import (
	"io"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/server"
)

type ServeRunner struct{}

func (r ServeRunner) Run(c api.Client, o GlobalOptions, w io.Writer) (err error) {
	err = server.Start()
	return
}

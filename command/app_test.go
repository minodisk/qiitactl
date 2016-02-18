package command_test

import (
	"os"
	"testing"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/command"
)

func TestNew(t *testing.T) {
	client := api.NewClient(func(subDomain, path string) (url string) {
		return
	})
	command.New(client, os.Stdout, os.Stderr)
}

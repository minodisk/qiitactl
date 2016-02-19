package command_test

import (
	"os"
	"testing"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/command"
	"github.com/minodisk/qiitactl/info"
	"github.com/minodisk/qiitactl/testutil"
)

var (
	inf = info.Info{
		Version: "0.0.0",
		TaskSettings: info.TaskSettings{
			GitHub: info.GitHub{
				Name: "qiitactl",
			},
		},
	}
)

func TestMain(m *testing.M) {
	code := m.Run()
	testutil.CleanUp()
	os.Exit(code)
}

func TestNew(t *testing.T) {
	client := api.NewClient(func(subDomain, path string) (url string) {
		return
	}, inf)
	command.New(inf, client, os.Stdout, os.Stderr)
}

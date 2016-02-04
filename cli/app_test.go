package cli_test

import (
	"testing"

	"github.com/minodisk/qiitactl/cli"
)

func TestNewApp(t *testing.T) {
	app := cli.NewApp()
	err := app.Run([]string{"qiitactl"})
	if err != nil {
		t.Error(err)
	}
}

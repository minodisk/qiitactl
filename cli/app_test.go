package cli_test

import (
	"testing"

	"github.com/minodisk/qiitactl/cli"
)

func TestNewApp(t *testing.T) {
	app := cli.GenerateApp()
	err := app.Run([]string{"qiitactl"})
	if err != nil {
		t.Error(err)
	}
}

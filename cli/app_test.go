package cli_test

import (
	"log"
	"os"
	"testing"

	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/cli"
)

func TestNewApp(t *testing.T) {
	err := os.Setenv("QIITA_ACCESS_TOKEN", "XXXXXXXXXXXX")
	if err != nil {
		log.Fatal(err)
	}
	client, err := api.NewClient(nil)
	if err != nil {
		t.Fatal(err)
	}
	app := cli.GenerateApp(client, os.Stdout)
	err = app.Run([]string{"qiitactl"})
	if err != nil {
		t.Error(err)
	}
}

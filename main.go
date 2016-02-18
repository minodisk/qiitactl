package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/command"
	"github.com/minodisk/qiitactl/info"
)

func main() {
	godotenv.Load()
	g := MustAsset(".goxc.json")
	_ = info.New(g)
	client := api.NewClient(nil)
	cmd := command.New(client, os.Stdout, os.Stderr)
	cmd.Run(os.Args)
}

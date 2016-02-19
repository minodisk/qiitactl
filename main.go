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
	info, err := info.New(g)
	if err != nil {
		panic("fail to load bindata")
	}
	client := api.NewClient(nil, info)
	cmd := command.New(info, client, os.Stdout, os.Stderr)
	cmd.Run(os.Args)
}

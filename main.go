package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/command"
)

func main() {
	godotenv.Load()
	client := api.NewClient(nil)
	cmd := command.New(client, os.Stdout, os.Stderr)
	cmd.Run(os.Args)
}

package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/cli"
)

func main() {
	godotenv.Load()
	client := api.NewClient(nil)
	app := cli.GenerateApp(client, os.Stdout, os.Stderr)
	app.Run(os.Args)
}

package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
	"github.com/minodisk/qiitactl/info"
)

func main() {
	godotenv.Load()

	app := cli.NewApp()
	app.Name = info.Name
	app.Version = info.Version
	app.Author = info.Author
	app.Usage = "Controls the Qiita posts"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}

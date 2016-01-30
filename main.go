package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	app := cli.NewApp()
	app.Name = "qiitactl"
	app.Version = Version
	app.Author = "minodisk"
	app.Usage = "Controls the Qiita posts"

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}

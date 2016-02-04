package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/minodisk/qiitactl/cli"
)

func main() {
	godotenv.Load()
	app := cli.NewApp()
	app.Run(os.Args)
}

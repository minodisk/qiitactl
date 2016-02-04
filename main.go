package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/minodisk/qiitactl/cli"
)

func main() {
	godotenv.Load()
	app := cli.GenerateApp()
	app.Run(os.Args)
}

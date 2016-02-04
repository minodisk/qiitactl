package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/minodisk/qiitactl/api"
	"github.com/minodisk/qiitactl/cli"
)

func main() {
	godotenv.Load()

	client, err := api.NewClient(nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	app := cli.GenerateApp(client, os.Stdout)
	app.Run(os.Args)
}

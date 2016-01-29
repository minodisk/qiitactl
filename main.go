package main

import (
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := cli.NewApp()
	app.Name = Name
	app.Version = Version
	app.Author = "minodisk"
	app.Email = ""
	app.Usage = ""

	app.Flags = GlobalFlags
	app.Commands = Commands
	app.CommandNotFound = CommandNotFound

	app.Run(os.Args)
}

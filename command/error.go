package command

import (
	"fmt"

	"github.com/codegangsta/cli"
)

func printError(c *cli.Context, err error) {
	if c.GlobalBool("debug") {
		panic(err)
	} else {
		fmt.Println(err)
	}
}

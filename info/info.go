package info

import "fmt"

type Info struct {
	Name        string
	Version     string
	Author      string
	Description string
}

func New(bindata []byte) (info Info) {
	fmt.Println(bindata)
	info.Name = "qiitactl"
	// info.Version =
	info.Author = "minodisk"
	info.Description = ""
	return
}

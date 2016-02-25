package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Message struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Data   interface{} `json:"data"`
}

func NewErrorRes(err error) Message {
	return Message{
		Method: "error",
		Data:   err.Error(),
	}
}

func NewReqRes(req Message, data interface{}) Message {
	return Message{
		ID:     req.ID,
		Method: req.Method,
		Data:   data,
	}
}

func Start() (err error) {
	http.Handle("/", http.FileServer(http.Dir("server/static")))
	http.Handle("/socket", websocket.Handler(socket))
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		return
	}
	return
}

func socket(ws *websocket.Conn) {
	for {
		var req Message
		websocket.JSON.Receive(ws, &req)
		switch req.Method {
		case "echo":
			io.Copy(ws, ws)
		case "ping":
			err := websocket.JSON.Send(ws, NewReqRes(req, "pong"))
			if err != nil {
				panic(err)
			}
		case "GetAllFiles":
			paths, err := findMarkdownFiles()
			if err != nil {
				websocket.JSON.Send(ws, NewErrorRes(err))
			}
			websocket.JSON.Send(ws, NewReqRes(req, paths))
		}
		time.Sleep(time.Second)
	}
}

type Element struct {
	ID       string     `json:"id"`
	Name     string     `json:"name"`
	Children []*Element `json:"children"`
}

func findMarkdownFiles() (root Element, err error) {
	dir, err := os.Getwd()
	if err != nil {
		return
	}
	dir = strings.TrimLeft(dir, "/")
	dir = strings.TrimRight(dir, "/")
	if dir == "" {
		dir = "."
	} else {
		dirs := strings.Split(dir, "/")
		dir = dirs[len(dirs)-1]
	}
	root.Name = dir

	err = filepath.Walk(".", func(path string, i os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
			return
		}
		if i.IsDir() {
			return
		}
		if filepath.Ext(path) != ".md" {
			return
		}

		parent := &root

		dir, file := filepath.Split(path)
		dir = strings.TrimLeft(dir, "/")
		dir = strings.TrimRight(dir, "/")
		if dir != "" {
			dirs := strings.Split(dir, "/")
			for _, name := range dirs {
				found := false
				for _, c := range parent.Children {
					if c.Name == name {
						parent = c
						found = true
						break
					}
				}

				if !found {
					e := &Element{Name: name}
					parent.Children = append(parent.Children, e)
					parent = e
				}
			}
		}

		parent.Children = append(parent.Children, &Element{Name: file})

		return
	})
	return
}

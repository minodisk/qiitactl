package server

import (
	"crypto/md5"
	"encoding/hex"
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
	Path     string     `json:"path"`
	Children []*Element `json:"children"`
}

func NewElement(path string, name string) (el Element) {
	el.Path = path
	hasher := md5.New()
	hasher.Write([]byte(el.Path))
	el.ID = hex.EncodeToString(hasher.Sum(nil))
	if name == "" {
		names := strings.Split(path, "/")
		el.Name = names[len(names)-1]
	} else {
		el.Name = name
	}
	return
}

func findMarkdownFiles() (root Element, err error) {
	path, err := os.Getwd()
	if err != nil {
		return
	}
	root = NewElement(path, "")

	var paths []string

	err = filepath.Walk(root.Path, func(path string, i os.FileInfo, e error) (err error) {
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

		path, err = filepath.Rel(root.Path, path)
		if err != nil {
			return
		}
		paths = append(paths, path)
		return
	})

	for _, path := range paths {
		names := strings.Split(path, "/")
		p := root.Path
		parent := &root
		for _, name := range names {
			p = filepath.Join(p, name)
			found := false
			var child *Element
			for _, child = range parent.Children {
				if child.Name == name {
					found = true
					break
				}
			}
			if !found {
				el := NewElement(p, name)
				child = &el
				parent.Children = append(parent.Children, child)
			}
			parent = child
		}
	}

	return
}

package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

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
	}
}

func findMarkdownFiles() (pathes []string, err error) {
	err = filepath.Walk(".", func(p string, i os.FileInfo, e error) (err error) {
		if e != nil {
			err = e
			return
		}
		if i.IsDir() {
			return
		}
		if filepath.Ext(p) != ".md" {
			return
		}
		pathes = append(pathes, p)
		return
	})
	return
}

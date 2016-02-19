package server

import (
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

func Start() (err error) {
	http.Handle("/", http.FileServer(http.Dir("./server/static")))
	http.Handle("/watcher", websocket.Handler(watcher))
	err = http.ListenAndServe(":9000", nil)
	if err != nil {
		return
	}
	return
}

func watcher(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

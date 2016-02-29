package server

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"

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
	router := httprouter.New()
	router.GET("/", handleIndex)
	router.GET("/markdown/*filepath", handleIndex)
	router.ServeFiles("/assets/*filepath", http.Dir("server/static/dist/assets"))
	router.Handler("GET", "/socket", websocket.Handler(socket))
	err = http.ListenAndServe(":9000", router)
	return
}

func handleIndex(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	b, err := ioutil.ReadFile("server/static/dist/index.html")
	if err != nil {
		panic(err)
	}
	w.Write(b)
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
	Rel      string     `json:"rel"`
	Abs      string     `json:"abs"`
	Children []*Element `json:"children"`
}

func NewElement(path, rootPath, name string) (el Element, err error) {
	el.Abs = path
	el.Rel, err = filepath.Rel(rootPath, path)
	if err != nil {
		return
	}
	hasher := md5.New()
	hasher.Write([]byte(el.Abs))
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
	rootPath, err := os.Getwd()
	if err != nil {
		return
	}
	root, err = NewElement(rootPath, rootPath, "")
	if err != nil {
		return
	}

	var paths []string

	err = filepath.Walk(rootPath, func(path string, i os.FileInfo, e error) (err error) {
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

		path, err = filepath.Rel(root.Abs, path)
		if err != nil {
			return
		}
		paths = append(paths, path)
		return
	})

	for _, path := range paths {
		names := strings.Split(path, "/")
		p := root.Abs
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
				el, err := NewElement(p, rootPath, name)
				if err != nil {
					return root, err
				}
				child = &el
				parent.Children = append(parent.Children, child)
			}
			parent = child
		}
	}

	return
}

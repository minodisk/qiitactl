package server

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"gopkg.in/fsnotify.v1"

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

func Start() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		err := initWatcher()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	go func() {
		err := initServer()
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	wg.Wait()
	return
}

var watcher *fsnotify.Watcher

func initWatcher() (err error) {
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	log.Println("init watcher")

	done := make(chan bool)
	log.Printf("num goroutune: %d", runtime.NumGoroutine())
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println(event.Name)
			case err := <-watcher.Errors:
				log.Println("error: ", err)
			}
		}
	}()
	<-done

	log.Println("init watcher complete")

	return
}

func watchFile(pathname string) (err error) {
	if err != nil {
		return
	}
	err = watcher.Add(pathname)
	return
}

func unwatchFile(pathname string) (err error) {
	err = watcher.Remove(pathname)
	return
}

func initServer() (err error) {
	log.Println("init server")

	router := httprouter.New()
	router.GET("/", handleIndex)
	router.GET("/markdown/*filepath", handleIndex)
	router.ServeFiles("/assets/*filepath", http.Dir("server/static/dist/assets"))
	router.Handler("GET", "/socket", websocket.Handler(socket))
	err = http.ListenAndServe(":9000", router)

	log.Println("init server complete")
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
	log.Println("start socket")

	wsChan := make(chan Message)

	go func() {
		log.Println("start receive")
		defer close(wsChan)
		var req Message
		err := websocket.JSON.Receive(ws, &req)
		if err != nil {
			log.Println(err)
			return
		}
		wsChan <- req
		log.Println("complete receive")
	}()

	for {
		// var req Message
		// websocket.JSON.Receive(ws, &req)
		select {
		case req, ok := <-wsChan:
			if !ok {
				log.Println("not ok", req)
				break
			}
			switch req.Method {
			case "GetAllFiles":
				paths, err := findMarkdownFiles()
				if err != nil {
					websocket.JSON.Send(ws, NewErrorRes(err))
				}
				websocket.JSON.Send(ws, NewReqRes(req, paths))
			case "WatchFile":
				err := watchFile(req.Data.(string))
				if err != nil {
					log.Println(err)
				}
			case "UnwatchFile":
				err := unwatchFile(req.Data.(string))
				if err != nil {
					log.Println(err)
				}
			}
		}
	}

	log.Println("complete socket")
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

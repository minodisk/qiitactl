package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Message struct {
	ID     string      `json:"id"`
	Method string      `json:"method"`
	Data   interface{} `json:"data"`
}

func NewErrorRes(id string, err error) Message {
	return Message{
		ID:     id,
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

func serveWebsocket(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		conn.Close()
	}()

	done := make(chan bool)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go reader(conn, wg, done)
	go writer(conn, wg, done)
	wg.Wait()
}

func reader(conn *websocket.Conn, wg *sync.WaitGroup, done chan bool) {
loop:
	for {
		select {
		case <-done:
			wg.Done()
			return
		default:
			var msg Message
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("read error: %s", err)
				break loop
			}

			log.Printf("receive: %s %v", msg.Method, msg.Data)

			switch msg.Method {
			case "GetAllFiles":
				paths, err := findMarkdownFiles()
				if err != nil {
					write(NewErrorRes(msg.ID, err))
					continue
				}
				write(NewReqRes(msg, paths))
			case "WatchFile":
				err := watchFile(msg.Data.(string))
				if err != nil {
					write(NewErrorRes(msg.ID, err))
				}
			case "UnwatchFile":
				err := unwatchFile(msg.Data.(string))
				if err != nil {
					write(NewErrorRes(msg.ID, err))
				}
			}
		}
	}

	close(done)
	wg.Done()
}

var writeChan = make(chan Message)

func writer(conn *websocket.Conn, wg *sync.WaitGroup, done chan bool) {
loop:
	for {
		select {
		case <-done:
			wg.Done()
			return
		case msg := <-writeChan:
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("write error: %s", err)
				break loop
			}
		}
	}

	close(done)
	wg.Done()
}

func write(msg Message) {
	writeChan <- msg
}

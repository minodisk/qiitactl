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

func NewErrorRes(req Message, err error) Message {
	return Message{
		ID:     req.ID,
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

	wg := new(sync.WaitGroup)
	wg.Add(2)
	go reader(conn, wg)
	go writer(conn, wg)
	wg.Wait()

	log.Println("end of handler")
}

func reader(conn *websocket.Conn, wg *sync.WaitGroup) {
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("read error: %s", err, msg)
			continue
		}

		log.Printf("receive: %s %v", msg.Method, msg.Data)

		switch msg.Method {
		case "GetAllFiles":
			paths, err := findMarkdownFiles()
			if err != nil {
				write(NewErrorRes(msg, err))
				continue
			}
			write(NewReqRes(msg, paths))
		case "WatchFile":
			err := watchFile(msg.Data.(string))
			if err != nil {
				write(NewErrorRes(msg, err))
			}
		case "UnwatchFile":
			err := unwatchFile(msg.Data.(string))
			if err != nil {
				write(NewErrorRes(msg, err))
			}
		}
	}
	wg.Done()
}

var writeChan = make(chan Message)

func writer(conn *websocket.Conn, wg *sync.WaitGroup) {
	for {
		select {
		case msg := <-writeChan:
			err := conn.WriteJSON(msg)
			if err != nil {
				log.Printf("write error: %s", err)
				break
			}
		}
	}
	wg.Done()
}

func write(msg Message) {
	writeChan <- msg
}

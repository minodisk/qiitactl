package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
)

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

func initServer() (err error) {
	log.Println("init server")

	router := httprouter.New()
	router.GET("/", handleIndex)
	router.GET("/markdown/*filepath", handleIndex)
	router.ServeFiles("/assets/*filepath", http.Dir("server/static/dist/assets"))
	router.GET("/socket", serveWebsocket)
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

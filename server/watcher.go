package server

import (
	"io/ioutil"
	"log"

	"gopkg.in/fsnotify.v1"
)

var watcher *fsnotify.Watcher

func initWatcher() (err error) {
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	log.Println("init watcher")

	// go func() {
	for {
		select {
		case event := <-watcher.Events:
			if (event.Op&fsnotify.Write == fsnotify.Write) ||
				(event.Op&fsnotify.Rename == fsnotify.Rename) {
				b, err := ioutil.ReadFile(event.Name)
				if err != nil {
					write(NewErrorRes("", err))
					continue
				}
				write(Message{Method: "ChangeFile", Data: File{Path: event.Name, Content: string(b)}})
			}

			watcher.Add(event.Name)
		case err := <-watcher.Errors:
			log.Println("error: ", err)
		}
	}
	// }()

	log.Println("init watcher complete")

	return
}

func watchFile(pathname string) (err error) {
	log.Printf("watchFile: %s", pathname)
	err = watcher.Add(pathname)
	return
}

func unwatchFile(pathname string) (err error) {
	log.Printf("unwatchFile: %s", pathname)
	err = watcher.Remove(pathname)
	return
}

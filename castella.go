package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

const (
	defaultPort = 8000
)

var (
	port = flag.Int("port", defaultPort, "server listen port")
	file = flag.String("file", "", "target file which watched update")
)

func main() {
	flag.Parse()
	addr := fmt.Sprintf(":%d", *port)

	hub := NewHub()
	go hub.run()

	watcher, err := NewWatcher(hub, *file)
	if err != nil {
		log.Println(err)
		return
	}
	go watcher.watch()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, watcher, w, r)
	})
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("listening port %d", *port)
}

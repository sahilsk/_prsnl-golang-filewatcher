package main

import (
	"filewatcher/lib"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
)

const (
	filename      = "index.html"
	watchfilename = "./testfile.log"
)

type tmplData struct {
	Title string
}

func main() {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ws", wsHandler)
	ipAndPort := "127.0.0.1" + ":" + "8080"
	listenErr := http.ListenAndServe(ipAndPort, nil)
	if listenErr != nil {
		log.Fatal("Error: ", listenErr)
	}
	fmt.Println("server started. Listening on " + ipAndPort)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

//HelloServer listen and serve requests
func wsHandler(w http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println(err)
		return
	}
	lib.FileWatcher(watchfilename, conn)
}

func rootHandler(w http.ResponseWriter, req *http.Request) {

	t := template.New("welcome")
	fread, readErr := ioutil.ReadFile(filename)
	if readErr != nil {
		log.Fatal("Error:", readErr)
	}
	tmpl, parseErr := t.Parse(string(fread))
	if parseErr != nil {
		log.Fatal("Error:", readErr)
	}
	tdata := tmplData{"FileWatcher"}
	err := tmpl.Execute(w, tdata)
	if err != nil {
		panic(err)
	}
}

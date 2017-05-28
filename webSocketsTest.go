package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/holwech/learnxyz-backend/models"
	_ "github.com/lib/pq"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Client subscribed")
	count := 1
	for {
		time.Sleep(2 * time.Second)
		err = conn.WriteMessage(websocket.TextMessage, count)
		if err != nil {
			fmt.Println(err)
			break
		}
		count++
	}
	fmt.Println("Client unsubscribed!")
}

func main() {
	index, err := ioutil.ReadAll(indexFile)
	if err != nil {
		fmt.Println(err)
	}
	http.HandleFunc("/ws", counterHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! :)")
	})
	fmt.Println("Serving at http://localhost:3000/topics/")
	http.ListenAndServe(":3000", nil)
}

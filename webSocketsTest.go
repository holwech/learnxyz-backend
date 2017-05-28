package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Counter struct {
	Count int `count:"int"`
}

// Sends a incrementing count to the client
func counterHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	counter := Counter{1}
	fmt.Println("Client subscribed")
	for {
		time.Sleep(2 * time.Second)
		fmt.Println("Sending ", counter)
		err = conn.WriteJSON(counter)
		fmt.Println("Message sent!")
		if err != nil {
			fmt.Println(err)
			break
		}
		counter.Count++
	}
	fmt.Println("Client unsubscribed!")
}

func main() {
	http.HandleFunc("/ws", counterHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello! :)")
	})
	fmt.Println("Serving websocket at ws://localhost:3000/ws/")
	http.ListenAndServe(":3000", nil)
}

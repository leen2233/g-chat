package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

type Conn struct {
	conn      *websocket.Conn
	nickname  string
}

var conns = []Conn{}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("New Connection")
	
	new_conn := Conn{
		conn: conn,
		nickname: getRandomNickname(),
	}
	conns = append(conns, new_conn)

	for _, v := range conns {
		v.conn.WriteMessage(websocket.TextMessage, fmt.Appendf([]byte(""), "[ %s joined the chat ]", new_conn.nickname))
	}

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		messageHandler(p, &new_conn)
	}
}


func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	log.Println("Starting server on 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

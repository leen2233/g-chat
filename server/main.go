package main

import (
	"encoding/json"
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

	msg := OutgoingMessage{
		Nickname: "system",
		Text: fmt.Sprintf("[%s joined the chat]", new_conn.nickname),
	}
	for _, c := range conns {
		sendMessage(c.conn, msg)
	}

	for {
		var jsonData []byte
		err := conn.ReadJSON(&jsonData)
		if err != nil {
			log.Println(err)
			return
		}

		msg := IncomingMessage{}
		json.Unmarshal(jsonData, &msg)
		messageHandler(msg, &new_conn)
	}
}


func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	log.Println("Starting server on 4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

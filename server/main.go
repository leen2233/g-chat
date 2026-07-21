package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
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

	msg := Message{
		Nickname: "system",
		Text: fmt.Sprintf("[%s joined the chat]", new_conn.nickname),
	}
	msgJson, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	e := Event{
		Type: "newMessage",
		Payload: msgJson,
	}
	for _, c := range conns {
		c.SendEvent(e)
	}
	
	
	new_conn.AddHandler("newMessage", newMessageHandler)


	new_conn.watchEvent()
}


func main(){
	host := flag.String("host", "127.0.0.1", "Host of server")
	port := flag.Int("port", 4000, "Port of server")

	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	address := fmt.Sprintf("%s:%s", *host, strconv.Itoa(*port))
	log.Printf("Starting server on %s\n", address)
	err := http.ListenAndServe(address, mux)
	log.Fatal(err)
}

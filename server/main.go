package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

var mappedConns map[int]*Conn

var conns = []*Conn{}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	new_conn := Conn{
		conn:     conn,
		Nickname: getRandomNickname(),
		Id:       getRandomId(),
	}
	conns = append(conns, &new_conn)
	mappedConns[new_conn.Id] = &new_conn

	log.Println("New Connection", new_conn.Nickname)

	// set identity on Connection
	identity := Identity{
		Nickname: new_conn.Nickname,
		Id: 			new_conn.Id,
	}
	sendEventHelper("setIdentity", identity, &new_conn)

	msg := ConnectedDisconnected{
		Nickname: new_conn.Nickname,
		Id: 			new_conn.Id,
		DateTime: time.Now(),
	}
	sendEventHelper("connected", msg, conns)	
	
	new_conn.AddHandler("newMessage", newMessageHandler)
	new_conn.AddHandler("getOnlineUsers", getOnlineUsersHandler)

	new_conn.watchEvent()
}


func main(){
	host := flag.String("host", "127.0.0.1", "Host of server")
	port := flag.Int("port", 4000, "Port of server")

	flag.Parse()

	mappedConns = make(map[int]*Conn)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	address := fmt.Sprintf("%s:%s", *host, strconv.Itoa(*port))
	log.Printf("Starting server on %s\n", address)
	err := http.ListenAndServe(address, mux)
	log.Fatal(err)
}


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
			if websocket.IsCloseError(err, 1006){
				// client disconnected
				out := OutgoingMessage{
					Nickname: "system",
					Text: fmt.Sprintf("[%s disconnected from chat]", new_conn.nickname),
				}
				for _, c := range conns {
					sendMessage(c.conn, out)
				}
			} else {
				log.Println(err)
			}
			return
		}

		msg := IncomingMessage{}
		json.Unmarshal(jsonData, &msg)
		messageHandler(msg, &new_conn)
	}
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

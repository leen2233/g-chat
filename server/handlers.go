package main

import (
	"encoding/json"
	"log"
	"time"
)


func newMessageHandler(payload any, conn *Conn) {
	msg, ok := payload.(*Message)
	if !ok {
		log.Println("payload is not *Message")
		return 
	}

	connTo, exists := mappedConns[msg.To]
	if exists {
		msg.From = conn.Id
		msg.DateTime = time.Now()

		sendEventHelper("newMessage", msg, connTo)
	}
}


func getOnlineUsersHandler(payload any, conn *Conn) {
	jsonData, err := json.Marshal(conns)
	if err != nil {
		log.Println(err)
	}

	e := Event{
		Type: "getOnlineUsers",
		Payload: jsonData,
	}

	conn.SendEvent(e)
}


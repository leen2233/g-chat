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

	msg.Nickname = conn.Nickname
	msg.DateTime = time.Now()

	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
	}
	e := Event{
		Type: "newMessage",
		Payload: jsonData,
	}
	for _, c := range conns {
		c.SendEvent(e)
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


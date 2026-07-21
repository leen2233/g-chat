package main

import (
	"encoding/json"
	"log"
)


func newMessageHandler(payload any, conn *Conn) {
	msg, ok := payload.(*Message)
	if !ok {
		log.Println("payload is not *Message")
		return 
	}

	msg.Nickname = conn.nickname

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


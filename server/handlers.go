package main

import (
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

		log.Println(conn.Id, conn.Nickname, msg.To, connTo.Id, connTo.Nickname)

		sendEventHelper("newMessage", msg, []*Conn{conn, connTo})
	}
}


func getOnlineUsersHandler(payload any, conn *Conn) {
	sendEventHelper("getOnlineUsers", conns, conn)
}


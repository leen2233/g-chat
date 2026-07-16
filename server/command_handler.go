package main

import (
	"fmt"
	"strings"

)

func messageHandler(msg IncomingMessage, conn *Conn){
	var message string
	// check if message is command
	if msg.Data[0] == byte('/') {
		command, args, _ := strings.Cut(msg.Data, " ")
		if command == "/set_nickname" {
			new_nickname := args
			old_nickname := conn.nickname
			conn.nickname = string(new_nickname)
			
			message = fmt.Sprintf("[ %s changed nickname to %s ]", old_nickname, conn.nickname)
		} else {
			message = "Unrecognized command"
		}
	} else {
		message = fmt.Sprintf("%s", msg.Data)
	}
	for _, c := range conns {
		msg := OutgoingMessage{
			Nickname: conn.nickname,
			Text: message,
		}
		sendMessage(c.conn, msg)
	}
}


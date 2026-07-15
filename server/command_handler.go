package main

import (
	"fmt"
	"strings"

	"github.com/gorilla/websocket"
)

func messageHandler(data []byte, conn *Conn){
	var message []byte
	// check if message is command
	if data[0] == byte('/') {
		command, args, _ := strings.Cut(string(data), " ")
		if command == "/set_nickname" {
			new_nickname := args
			old_nickname := conn.nickname
			conn.nickname = string(new_nickname)
			
			message = fmt.Appendf([]byte(""), "[ %s changed nickname to %s ]", old_nickname, conn.nickname)
		} else {
			message = []byte("Unrecognized command")
		}
	} else {
		message = fmt.Appendf([]byte(""), "[%s] %s", conn.nickname, data)
	}
	for _, c := range conns {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}


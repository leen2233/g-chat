package main

import (
	"encoding/json"

	"github.com/gorilla/websocket"
)

func sendMessage(conn *websocket.Conn, message OutgoingMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	conn.WriteJSON(data)

	return nil
}


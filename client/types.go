package main

import "encoding/json"


type Event struct {
	Type     		string					`json:"type"`
	Payload  		json.RawMessage `json:"payload"`
}


type Message struct {
	Nickname    string
	Text        string
}

func newMessage() any {
	return &Message{}
}


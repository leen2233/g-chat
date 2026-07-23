package main

import (
	"encoding/json"
	"time"
)


type Event struct {
	Type     		string					`json:"type"`
	Payload  		json.RawMessage `json:"payload"`
}


type Message struct {
	From		    int				`json:"from"`
	To  				int 			`json:"to"`
	Text        string    `json:"text"`
	DateTime    time.Time `json:"datetime"`
}


type ConnectedDisconnected struct {
	Nickname    string    `json:"nickname"`
	Id          int 			`json:"id"`
	DateTime    time.Time `json:"datetime"`
}


type Identity struct {
	Nickname 		string 		`json:"nickname"`
	Id 					int 			`json:"id"`
}


func newMessage() any {
	return &Message{}
}


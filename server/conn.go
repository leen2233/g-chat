package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)


var EventPayloadMap = map[string]func() any{
	"newMessage": newMessage,
}

type Conn struct {
	Handlers     map[string]func(payload any, conn *Conn)
	conn       *websocket.Conn
	nickname  string
}

func (c *Conn) watchEvent() {
	for {
		var jsonData []byte
		err := c.conn.ReadJSON(&jsonData)
		if err != nil {
			if websocket.IsCloseError(err, 1006){
				// client disconnected
				msg := Message{
					Nickname: "system",
					Text: fmt.Sprintf("[%s disconnected from chat]", c.nickname),
				}
				msgJson, err := json.Marshal(msg)
				if err != nil {
					log.Println(err)
				}
				e := Event{
					Type: "newMessage",
					Payload: msgJson,
				}
				for _, c := range conns {
					c.SendEvent(e)
				}
			} else {
				log.Println(err)
			}
			return
		}
		
		e := Event{}
		err = json.Unmarshal(jsonData, &e)
		if err != nil {
			log.Println(err)
		}
		c.HandleEvent(e)
	}

}

func (c *Conn) HandleEvent(e Event) error {
	var payload any
	constructor := EventPayloadMap[e.Type]
	if constructor != nil {
		payload = constructor()
		err := json.Unmarshal(e.Payload, payload)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unknown constructor for event type: (%s). Set it to use it", e.Type)
	}

	handler := c.Handlers[e.Type]
	handler(payload, c)

	return nil
}

func (c *Conn) SendEvent(e Event) error {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return err
	}

	c.conn.WriteJSON(jsonData)
	return nil
}


func (c *Conn) AddHandler(eventType string, callback func(payload any, conn *Conn)) {
	if c.Handlers == nil {
		c.Handlers = make(map[string]func(payload any, conn *Conn))
	}
	c.Handlers[eventType] = callback
}


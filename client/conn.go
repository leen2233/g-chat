package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)


var EventPayloadMap = map[string]func() any{
	"newMessage": newMessage,
}

type Conn struct {
	Host         string
	Port         int
	Handlers     map[string][]func(payload any)
	WsConn       *websocket.Conn
}

func (c *Conn) Connect() error {
	address := fmt.Sprintf("ws://%s:%s", c.Host, strconv.Itoa(c.Port))

	dialer := websocket.Dialer{}
	ws_conn, _, err := dialer.Dial(address, nil)
	c.WsConn = ws_conn
	if err != nil {
		return err
	}

	go func() {
		for {
			var jsonData []byte
			err := ws_conn.ReadJSON(&jsonData)
			if err != nil {
				log.Fatal(err)
			}

			e := Event{}
			err = json.Unmarshal(jsonData, &e)	
			if err != nil {
				log.Print(err)
			}
			err = c.HandleEvent(e)
			if err != nil {
				log.Print(err)
			}
		}
	}()

	return nil
}


func (c *Conn) HandleEvent(e Event) error {
	var payload any
	constructor := EventPayloadMap[e.Type]
	if constructor != nil {
		payload = constructor()
		err = json.Unmarshal(e.Payload, payload)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("Unknown constructor for event type: (%s). Set it to use it", e.Type)
	}

	handlers := c.Handlers[e.Type]
	for _, handle := range handlers {
		handle(payload)
	}

	return nil
}

func (c *Conn) SendEvent(e Event) error {
	jsonData, err := json.Marshal(e)
	if err != nil {
		return err
	}

	c.WsConn.WriteJSON(jsonData)
	return nil
}


func (c *Conn) AddHandler(eventType string, callback func(payload any)) {
	if c.Handlers == nil {
		c.Handlers = make(map[string][]func(payload any))
	}
	c.Handlers[eventType] = append(c.Handlers[eventType], callback)
}


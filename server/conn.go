package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)


var EventPayloadMap = map[string]func() any{
	"newMessage": newMessage,
}

type Conn struct {
	handlers    map[string]func(payload any, conn *Conn)
	conn        *websocket.Conn
	Nickname  	string		`json:"nickname"`
	Id         	int 			`json:"id"`
}

func (c *Conn) watchEvent() {
	for {
		var jsonData []byte
		err := c.conn.ReadJSON(&jsonData)
		if err != nil {

			// if websocket.IsCloseError(err, 1006){
				// client disconnected
				msg := ConnectedDisconnected{
					Nickname: c.Nickname,
					Id:				c.Id,
					DateTime: time.Now(),
				}
				msgJson, err := json.Marshal(msg)
				if err != nil {
					log.Println(err)
				}
				e := Event{
					Type: "disconnected",
					Payload: msgJson,
				}
				indexOfConn := -1
				for index, v := range conns {
					v.SendEvent(e)
					if c.Id == v.Id {
						// try to find index of disconnected connection to later remove it
						indexOfConn = index
					}
				}
				if indexOfConn >= 0 {
					// if index found, remove disconnected connection from conns.
					conns = append(conns[:indexOfConn], conns[indexOfConn+1:]...)
					// TODO: remove conn from mappedConns 
				}
				log.Println(indexOfConn, conns)
			// } else {
			// 	log.Println(err)
			// }
			return
		}
		
		e := Event{}
		err = json.Unmarshal(jsonData, &e)
		if err != nil {
			log.Println(err)
		}
		err = c.HandleEvent(e)
		if err != nil {
			log.Println(err)
		}
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
	}

	log.Println(fmt.Sprintf("[New event] %s   [payload]: %+v", e.Type, payload))

	handler := c.handlers[e.Type]
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
	if c.handlers == nil {
		c.handlers = make(map[string]func(payload any, conn *Conn))
	}
	c.handlers[eventType] = callback
}


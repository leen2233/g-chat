package main

import (
	"encoding/json"
	"log"
)


func sendEventHelper(eventType string, v any, target any) {
	vJson, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}
	e := Event{
		Type: eventType,
		Payload: vJson,
	}

	switch conns := target.(type) {
	case []*Conn:
		for _, conn := range conns {
			conn.SendEvent(e)
		}
	case *Conn:
		conns.SendEvent(e)
	}

	log.Printf("[SENT] %s , payload: %+v , conn: %+v\n", e.Type, v, target)
}


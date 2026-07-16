package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jroimartin/gocui"
)

var messagesView *gocui.View
var g *gocui.Gui
var incoming = make(chan IncomingMessage, 100)

var conn *websocket.Conn
var err error

func main() {
	dialer := websocket.Dialer{}
	conn, _, err = dialer.Dial("ws://127.0.0.1:4000", nil)
	if err != nil {
		log.Fatal("Couldn't connect to server")
	}
	go func() {
		for {
			var jsonData []byte
			err := conn.ReadJSON(&jsonData)
			if err != nil {
				log.Fatal(err)
			}
			
			msg := IncomingMessage{}
			err = json.Unmarshal(jsonData, &msg)
			if err != nil {
				log.Println(err)
			}
			incoming <- msg
		}
	}()

	// GUI
	g, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}


func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	messagesView, err = g.SetView("messages", 1, 1, maxX - 2, maxY - 5)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	messagesView.Autoscroll = true

	go func (){
		for msg := range incoming {
			g.Update(func(g *gocui.Gui) error {
					fmt.Fprintf(messagesView, "\x1b[0;32m%s \x1b[0;0m %s\n", msg.Nickname, msg.Text)
					return nil
				},
			)
		}
	}()

	input, err := g.SetView("input", 1, maxY - 4, maxX - 2, maxY - 2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	input.Editable = true
	input.Editor = Editor
	g.Cursor = true
	g.SetCurrentView("input")	

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}


var Editor gocui.Editor = gocui.EditorFunc(Edit)

func Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
		case key == gocui.KeyEnter:
			message := strings.TrimSpace(v.Buffer())

			if conn != nil && message != "" {
				msg := OutgoingMessage{
					Data: message,
				}
				jsonData, err := json.Marshal(msg)
				if err != nil {
					fmt.Println(err)
				}
				conn.WriteJSON(jsonData)
				v.Clear()
				v.SetCursor(0, 0)
			} 
		case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
			v.EditDelete(true)
		case key == gocui.KeySpace:
			v.EditWrite(' ')
		case ch != 0 && mod == 0:
			v.EditWrite(ch)
	}
}



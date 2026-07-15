package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/jroimartin/gocui"
)

var messagesView *gocui.View
var g *gocui.Gui
var unreadMessages []string

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
			_, p, err := conn.ReadMessage()
			if err != nil {
				log.Fatal(err)
			}

			if messagesView != nil {
				fmt.Fprintf(messagesView, "%s\n", string(p))
				g.Update(func(g *gocui.Gui) error {return nil})
			} else {
				unreadMessages = append(unreadMessages, string(p))
			}
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
	for _, v := range unreadMessages {
		fmt.Fprintf(messagesView, "%s\n", v)
	}
	unreadMessages = nil

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
	if key == gocui.KeyEnter {
		message := strings.TrimSpace(v.Buffer())

		if conn != nil {
			conn.WriteMessage(websocket.TextMessage, []byte(message))
			v.Clear()
			v.SetCursor(0, 0)
		}
	} else if key == gocui.KeyBackspace || key == gocui.KeyBackspace2 {
		v.EditDelete(true)
	} else {
		v.EditWrite(ch)
	}
}



package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"
)

var messagesView *gocui.View
var g *gocui.Gui
var incoming = make(chan *Message, 100)

var conn *Conn
var err error

func main() {
	host := flag.String("host", "127.0.0.1", "Host of server")
	port := flag.Int("port", 4000, "Port of server")

	flag.Parse()

	conn = &Conn{
		Host: *host,
		Port: *port,
	}
	err = conn.Connect()
	if err != nil {
		log.Fatal(err)
	}

	conn.AddHandler("newMessage", func(payload any) {
		msg, ok := payload.(*Message)
		if !ok {
			log.Println("payload is not *Message")
			return
		}
		incoming <- msg
	})

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
	chatsView, err := g.SetView("chats", 1, 1, 28, maxY - 5)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	chatsView.Frame = true
	chatsView.Title = "Chats"

	messagesView, err = g.SetView("messages", 30, 1, maxX - 2, maxY - 5)
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
				msg := Message{
					Text: message,
					Nickname: "",
				}
				jsonMsg, err := json.Marshal(msg)
				if err != nil {
					log.Print(err)
				}
				e := Event{
					Type: "newMessage",
					Payload: jsonMsg,
				}
				err = conn.SendEvent(e)
				if err != nil {
					log.Print(err)
				}

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



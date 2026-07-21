package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


var conn *Conn
var err error
var input *tview.TextArea
var messagesBox *tview.TextView

func main() {
	host := flag.String("host", "127.0.0.1", "Host of server")
	port := flag.Int("port", 4000, "Port of server")

	flag.Parse()

	chatsBox := tview.NewBox().SetBorder(true).SetTitle("Chats")
	messagesBox = tview.NewTextView().SetDynamicColors(true)
	messagesBox.SetBorder(true)
	input = tview.NewTextArea().SetPlaceholder("type a message...").SetPlaceholderStyle(inputPlaceholderStyle)
	input.SetBorder(true)
	input.SetChangedFunc(inputChanged)
	input.SetInputCapture(inputKeyHandler)

	mainGrid := tview.NewGrid().SetColumns(30, 0).SetRows(0, 4)
	mainGrid.AddItem(chatsBox, 0, 0, 2, 1, 0, 30, false)
	mainGrid.AddItem(messagesBox, 0, 1, 1, 1, 0, 100, false)
	mainGrid.AddItem(input, 1, 1, 1, 1, 0, 100, true)


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
		
		fmt.Fprintf(messagesBox, "[green]%s[-] [gray]15:04[-]\n%s\n\n", msg.Nickname, msg.Text)
	})
	
	if err := tview.NewApplication().SetRoot(mainGrid, true).Run(); err != nil {
		panic(err)
	}
}


func inputChanged() {
}

func inputKeyHandler(e *tcell.EventKey) *tcell.EventKey {
	if e.Key() == tcell.KeyEnter {
		if e.Modifiers() == 4 {
			// shift + enter pressed
			return e
		} else {
			// send a message
			message := Message{
				Text: input.GetText(),
			}
			jsonData, err := json.Marshal(message)
			if err != nil {
				log.Println("couldn't marshal json'")
			}
			e := Event{
				Type: "newMessage",
				Payload: jsonData,
			}
			conn.SendEvent(e)

			input.SetText("", true)
			return nil
		}
	}
	return e
}


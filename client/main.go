package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)


var conn 					*Conn
var err 					error
var input 				*tview.TextArea
var messagesBox 	*tview.TextView
var chatsBox    	*tview.List
var app 					*tview.Application

var chatsList     map[int]*OnlineUser
var selectedChat  *OnlineUser


func main() {
	host := flag.String("host", "127.0.0.1", "Host of server")
	port := flag.Int("port", 4000, "Port of server")

	flag.Parse()

	chatsBox = tview.NewList()
	chatsBox.ShowSecondaryText(false).SetBorder(true).SetTitle("Chats")
	chatsBox.SetChangedFunc(chatListChangedHandler)
	messagesBox = tview.NewTextView().SetDynamicColors(true)
	messagesBox.SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	input = tview.NewTextArea().SetPlaceholder("type a message...").SetPlaceholderStyle(inputPlaceholderStyle)
	input.SetBorder(true)
	input.SetInputCapture(inputKeyHandler)

	mainGrid := tview.NewGrid().SetColumns(30, 0).SetRows(0, 4)
	mainGrid.AddItem(chatsBox, 0, 0, 2, 1, 0, 30, false)
	mainGrid.AddItem(messagesBox, 0, 1, 1, 1, 0, 100, false)
	mainGrid.AddItem(input, 1, 1, 1, 1, 0, 100, true)

	chatsList = make(map[int]*OnlineUser)

	conn = &Conn{
		Host: *host,
		Port: *port,
	}
	err = conn.Connect()
	if err != nil {
		log.Fatal(err)
	}

	conn.AddHandler("newMessage", handleNewMessage)
	conn.AddHandler("connected", handleConnected)
	conn.AddHandler("disconnected", handleDisconnected)
	conn.AddHandler("getOnlineUsers", handleGetOnlineUsers)
	conn.AddHandler("setIdentity", handleSetIdentity)
	
	app = tview.NewApplication().SetRoot(mainGrid, true)
	app.EnableMouse(true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}

// enter handling
func inputKeyHandler(e *tcell.EventKey) *tcell.EventKey {
	if e.Key() == tcell.KeyEnter {
		if e.Modifiers() == 4 {
			// shift + enter pressed
			return e
		} else {
			// send a message
			if strings.TrimSpace(input.GetText()) == "" {
				// don't send a message if it's empty
				return nil
			}
			
			message := Message{
				Text: input.GetText(),
				To:		selectedChat.Id,
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


func chatListChangedHandler(index int, main string, secondary string, shortcut rune) {
	// handle when chatlist value changed. this function is fired when user selects or just moving between value without selecting or new item added and marked as selected.
	if chatsBox.GetItemCount() > 0 {
		// try to get current selected value if count is > 0
		_, secondary := chatsBox.GetItemText(chatsBox.GetCurrentItem())
		selectedChatId, err := strconv.Atoi(secondary)
		if err != nil {
			log.Println("couldn't convert selectedChatId to int")
			return
		}
		if selectedChat == nil || selectedChatId != selectedChat.Id {
			// new selected chat id is not same with previously selected chat
			// load messages of newly selected chat
			selectedChat = chatsList[selectedChatId]
			messagesBox.Clear()
			for _, message := range selectedChat.Messages {
				var nickname string
				if conn.Id == message.To {
					nickname = chatsList[message.From].Nickname
				} else {
					nickname = conn.Nickname
				}

				fmt.Fprintf(messagesBox, "[green]%s[-] [gray]%s[-]\n%s\n\n", nickname, message.DateTime.Format("15:04:05"), message.Text)
			}
			messagesBox.ScrollToEnd()
			if app != nil {
				app.Draw()
			}
		}
	}


}



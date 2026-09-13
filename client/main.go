package main

import (
	"flag"
	"fmt"

	// "fmt"
	"log"
	"strconv"

	"github.com/rivo/tview"
)


var conn 					*Conn
var err 					error
var input 				*tview.TextArea
var messagesBox 	*tview.TextView
var identityBox   *tview.TextView
var chatsBox    	*tview.List
var app 					*tview.Application

var chatsList     map[int]*OnlineUser
var selectedChat  *OnlineUser


func main() {
	host := flag.String("host", "127.0.0.1", "Host of server")
	port := flag.Int("port", 4000, "Port of server")

	flag.Parse()

	chatsBox = tview.NewList()
	chatsBox.ShowSecondaryText(true).SetBorder(true).SetTitle("Chats").SetBorderPadding(0, 0, 1, 2)
	chatsBox.SetChangedFunc(chatListChangedHandler)
	messagesBox = tview.NewTextView().SetDynamicColors(true)
	messagesBox.SetBorder(true)
	identityBox = tview.NewTextView().SetScrollable(false)
	identityBox.SetBorder(true).SetBorderPadding(0, 0, 1, 1)
	input = tview.NewTextArea().SetPlaceholder("type a message...").SetPlaceholderStyle(inputPlaceholderStyle)
	input.SetBorder(true)
	input.SetInputCapture(inputKeyHandler)

	mainGrid := tview.NewGrid().SetColumns(30, 0).SetRows(0, 4)
	mainGrid.AddItem(chatsBox, 0, 0, 1, 1, 0, 30, false)
	mainGrid.AddItem(identityBox, 1, 0, 1, 1, 0, 30, false)
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
	app.SetInputCapture(globalInputCapture)
	if err := app.Run(); err != nil {
		panic(err)
	}
}


func chatListChangedHandler(index int, main string, secondary string, shortcut rune) {
	// handle when chatlist value changed. this function is fired when user selects or just moving between value without selecting or new item added and marked as selected.
	if chatsBox.GetItemCount() > 0 {
		selectedChatId, err := strconv.Atoi(secondary)
		if err != nil {
			log.Panic("couldn't convert selectedChatId to int")
			return
		}

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
	}
}



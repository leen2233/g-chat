package main

import (
	"fmt"
	"log"
	"strconv"
)


func handleNewMessage(payload any) {
	msg, ok := payload.(*Message)
	if !ok {
		log.Println("payload is not *Message")
		return
	}
	
	fmt.Fprintf(messagesBox, "[green]%s[-] [gray]%s[-]\n%s\n\n", msg.Nickname, msg.DateTime.Format("15:04:05"), msg.Text)
	messagesBox.ScrollToEnd()
	if app != nil {
		app.Draw()
	}
}

func handleConnected(payload any) {
	msg, ok := payload.(*ConnectedDisconnected)
	if !ok {
		log.Println("payload is not *ConnectedDisconnected")
		return
	}

	fmt.Fprintf(messagesBox, "[green]%s connected to chat[-] [gray]%s[-]\n\n", msg.Nickname, msg.DateTime.Format("15:04:05"))

	// search and add connection to chats list if not exists
	if _, exists := chatsList[msg.Id]; !exists {
		chatsBox.AddItem(msg.Nickname, strconv.Itoa(msg.Id), 'a', nil)
		chatsList[msg.Id] = OnlineUser{
			Nickname: msg.Nickname,
			Id: msg.Id,
		}
	}

	if app != nil {
		app.Draw()
	}
}

func handleDisconnected(payload any) {
	msg, ok := payload.(*ConnectedDisconnected)
	if !ok {
		log.Println("payload is not *ConnectedDisconnected")
		return
	}

	fmt.Fprintf(messagesBox, "[red]%s disconnected from chat[-] [gray]%s[-]\n\n", msg.Nickname, msg.DateTime.Format("15:04:05"))
	
	// remove this connection from chats list
	for i := 0 ; i < chatsBox.GetItemCount(); i++ {
		if _, secondary := chatsBox.GetItemText(i); secondary == strconv.Itoa(msg.Id) {
			chatsBox.RemoveItem(i)
		}
	}

	if app != nil {
		app.Draw()
	}
}


func handleGetOnlineUsers(payload any) {
	users, ok := payload.(*[]OnlineUser)
	if !ok {
		log.Println("payload is not *[]OnlineUsers")
	}

	for _, user := range *users {
		if _, exists := chatsList[user.Id]; !exists {
			chatsBox.AddItem(user.Nickname, strconv.Itoa(user.Id), 'a', func(){})
			chatsList[user.Id] = user
		}
	}
}


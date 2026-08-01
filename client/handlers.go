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

	// messages should be contained in peer user whether it's send by self or another
	var peerId  	int
	var nickname 	string
	if conn.Id == msg.To {
		peerId = msg.From
		nickname = chatsList[msg.From].Nickname
	} else {
		peerId = msg.To
		nickname = conn.Nickname
	}

	if chatsList[peerId] == nil {
		chatsList[peerId] = &OnlineUser{
			Id: peerId,
		}
	}
	chatsList[peerId].Messages = append(chatsList[peerId].Messages, msg)
	
	// add message to message box if peer chat is selected
	if peerId == selectedChat.Id {
		fmt.Fprintf(messagesBox, "[green]%s[-] [gray]%s[-]\n%s\n\n", nickname, msg.DateTime.Format("15:04:05"), msg.Text)
		messagesBox.ScrollToEnd()
		if app != nil {
			app.Draw()
		}
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
	if msg.Id != conn.Id {
		if _, exists := chatsList[msg.Id]; !exists {
			chatsList[msg.Id] = &OnlineUser{
				Nickname: msg.Nickname,
				Id: msg.Id,
			}
			chatsBox.AddItem(msg.Nickname, strconv.Itoa(msg.Id), 'a', nil)
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

	chatsList[msg.Id]	= nil
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
		if user.Id != conn.Id {
			if _, exists := chatsList[user.Id]; !exists {
				chatsList[user.Id] = &user
				chatsBox.AddItem(user.Nickname, strconv.Itoa(user.Id), 'a', func(){})
			}
		}
	}
}


func handleSetIdentity(payload any) {
	identity, ok := payload.(*Identity)
	if !ok {
		log.Println("payload is not *Identity")
	}
	
	conn.Nickname = identity.Nickname
	conn.Id				= identity.Id

	for i := 0 ; i < chatsBox.GetItemCount(); i++ {
		if _, secondary := chatsBox.GetItemText(i); secondary == strconv.Itoa(conn.Id) {
			chatsBox.RemoveItem(i)
		}
	}
}


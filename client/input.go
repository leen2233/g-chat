package main

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
)

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
			if selectedChat == nil {
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


func globalInputCapture(e *tcell.EventKey) *tcell.EventKey {
	if e.Key() == tcell.KeyEsc {
		app.SetFocus(nil)
	} else {
		// check focused element. this allows setting keybinds for specific boxes
		if app.GetFocus() == nil {  // no focused element
			if e.Rune() == 99 {
				// if pressed character is "c", then focus chats list
				app.SetFocus(chatsBox)
			} else if e.Rune() == 109 {
				// if pressed character is "m", then focus messages list
				app.SetFocus(messagesBox)
			} else if e.Rune() == 13 {
				// if pressed character is "Enter", then focus message input
				app.SetFocus(input)
			}
		} else if app.GetFocus() == chatsBox {
			// allow user to move with j/k 
			if e.Rune() == 106 {
				chatsBox.SetCurrentItem(chatsBox.GetCurrentItem() + 1)
			} else if e.Rune() == 107 {
				chatsBox.SetCurrentItem(chatsBox.GetCurrentItem() - 1)
			} else if e.Rune() == 13 {
				// if pressed character is "Enter", then focus message input
				app.SetFocus(input)
			}
		}

	}

	return e
}


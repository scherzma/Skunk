package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// This Peer gets invited to a chat
type InviteToChatHandler struct {
	userChatLogic chat.ChatLogic
}

func (i *InviteToChatHandler) HandleMessage(message network.Message) error {

	// Structure of the message:
	/*
		{
			"chatId": "asdf",
			"chatName": "asdf",
			"peers": [
				"asdf",
				"asdf"
			]
		}
	*/

	var content struct {
		ChatID   string `json:"chatId"`
		ChatName string `json:"chatName"`
		Peers    []string
	}

	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	i.userChatLogic.ReceiveChatInvitation(message.FromUser, content.ChatID, content.ChatName, content.Peers)
	//TODO store received chat invitation for later use
	// For example: if the user wants to join the chat

	return nil
}

package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type InviteToChatHandler struct {
	userChatLogic chat.ChatLogic
}

func (i *InviteToChatHandler) HandleMessage(message network.Message) error {

	// chatRepo := p_model.GetNetworkChatsInstance()

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

	i.userChatLogic.ProcessMessageForUser(message)

	return nil
}

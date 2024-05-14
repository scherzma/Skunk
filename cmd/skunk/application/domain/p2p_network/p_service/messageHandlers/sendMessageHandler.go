package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type SendMessageHandler struct {
	userChatLogic chat.ChatLogic
}

func (s *SendMessageHandler) HandleMessage(message network.Message) error {

	// Structure of the message:
	/*
		{
			"message": "asdfasdfasdf",
		}
	*/

	var content struct {
		Message string `json:"message"`
	}

	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	s.userChatLogic.RecieveMessage(message.SenderID, message.ChatID, content.Message)

	return nil
}

package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type SetUsernameHandler struct {
	userChatLogic chat.ChatLogic
}

func (s *SetUsernameHandler) HandleMessage(message network.Message) error {

	// Structure of the message:
	/*
		{
			"message": "asdfasdfasdf",
		}
	*/

	var content struct {
		Username string `json:"username"`
	}

	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	s.userChatLogic.PeerSetsUsername(message.SenderID, message.ChatID, content.Username)

	return nil
}

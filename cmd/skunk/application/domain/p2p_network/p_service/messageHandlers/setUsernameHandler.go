package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
)

type setUsernameHandler struct {
	userChatLogic      chat.ChatLogic
	userMessageStorage store.UserMessageStoragePort
}

func NewSetUsernameHandler(userChatLogic chat.ChatLogic, userMessageStorage store.UserMessageStoragePort) *setUsernameHandler {
	return &setUsernameHandler{
		userChatLogic:      userChatLogic,
		userMessageStorage: userMessageStorage,
	}
}

func (s *setUsernameHandler) HandleMessage(message network.Message) error {

	// Structure of the message:
	/*
		{
			"username": "asdfasdfasdf",
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

	// Store the username
	err = s.userMessageStorage.PeerSetUsername(message.SenderID, message.ChatID, content.Username)
	if err != nil {
		fmt.Println("Error storing username")
		return err
	}

	// Handle the username change
	s.userChatLogic.PeerSetsUsername(message.SenderID, message.ChatID, content.Username)

	return nil
}

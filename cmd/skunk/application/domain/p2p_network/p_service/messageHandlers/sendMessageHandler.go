package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
)

type sendMessageHandler struct {
	userChatLogic         chat.ChatLogic
	networkMessageStorage store.NetworkMessageStoragePort
}

func NewSendMessageHandler(userChatLogic chat.ChatLogic, networkMessageStorage store.NetworkMessageStoragePort) *sendMessageHandler {
	return &sendMessageHandler{
		userChatLogic:         userChatLogic,
		networkMessageStorage: networkMessageStorage,
	}
}

func (s *sendMessageHandler) HandleMessage(message network.Message) error {

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

	// Store the message
	err = s.networkMessageStorage.StoreMessage(message)
	if err != nil {
		fmt.Println("Error storing message")
		return err
	}

	// Handle the received message
	s.userChatLogic.ReceiveMessage(message.SenderID, message.ChatID, content.Message)

	return nil
}

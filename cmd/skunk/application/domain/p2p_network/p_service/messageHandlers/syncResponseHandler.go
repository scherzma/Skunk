package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
)

type syncResponseHandler struct {
	networkMessageStorage store.NetworkMessageStoragePort
}

func NewSyncResponseHandler(networkMessageStorage store.NetworkMessageStoragePort) *syncResponseHandler {
	return &syncResponseHandler{
		networkMessageStorage: networkMessageStorage,
	}
}

// HandleMessage processes the received "SyncResponse" message.
// It retrieves the chat message repository, unmarshals the received messages from the message content,
// and adds each received message to the chat message repository.
// TODO: Implement a security check to ensure that the message is valid.
func (s *syncResponseHandler) HandleMessage(message network.Message) error {
	//chatRepo := p_model.GetNetworkChatsInstance() TODO: change
	//chatMessageRepo := chatRepo.GetChat(message.ChatID) TODO: change

	var receivedMessages []network.Message
	err := json.Unmarshal([]byte(message.Content), &receivedMessages)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	for _, msg := range receivedMessages {
		// Store the message
		err = s.networkMessageStorage.StoreMessage(msg)
		if err != nil {
			fmt.Println("Error storing message:", err)
			return err
		}
		// Add the message to the chat repository
		//chatMessageRepo.AddMessage(msg) TODO: change
	}

	return nil
}

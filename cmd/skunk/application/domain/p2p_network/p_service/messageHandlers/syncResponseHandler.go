package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// SyncResponseHandler handles the "SyncResponse" message operation.
type SyncResponseHandler struct {
}

// HandleMessage processes the received "SyncResponse" message.
// It retrieves the chat message repository, unmarshals the received messages from the message content,
// and adds each received message to the chat message repository.
// TODO: Implement a security check to ensure that the message is valid.
func (s *SyncResponseHandler) HandleMessage(message network.Message) error {

	chatRepo := p_model.GetNetworkChatsInstance()
	chatMessageRepo := chatRepo.GetChat(message.ChatID)

	var receivedMessages []network.Message
	err := json.Unmarshal([]byte(message.Content), &receivedMessages)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	for _, message := range receivedMessages {
		chatMessageRepo.AddMessage(message)
	}

	return nil
}

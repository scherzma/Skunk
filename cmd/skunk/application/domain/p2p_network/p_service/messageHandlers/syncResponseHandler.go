package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type SyncResponseHandler struct {
}

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

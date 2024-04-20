package messageHandlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"time"
)

type SyncRequestHandler struct {
}

func (s *SyncRequestHandler) HandleMessage(message network.Message) error {

	chatRepo := p_model.GetNetworkChatsInstance()
	chatMessageRepo := chatRepo.GetChat(message.ChatID)

	// Parse the content of the message
	/*
		{
		  "existingMessageIds": [
			"<message id 1>",
			"<message id 2>",
			...
		  ]
		}
	*/
	var content []string
	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		return err
	}

	// find difference between "message" already known messages and own messages that the other peer does not know
	missingExternalMessages := chatMessageRepo.GetMissingExternalMessages(content)
	missingInternalMessages := chatMessageRepo.GetMissingInternalMessages(content)

	// Convert missingExternalMessages to a JSON string
	externalMessagesBytes, err := json.Marshal(missingExternalMessages)
	if err != nil {
		return err
	}

	// Convert missingInternalMessages to a JSON string
	internalMessagesBytes, err := json.Marshal(missingInternalMessages)
	if err != nil {
		return err
	}

	//TODO send the sync response to the other peer
	syncResponse := network.Message{
		Id:        uuid.New().String(),
		Timestamp: time.Now().UnixNano(),
		Content:   string(externalMessagesBytes),
		FromUser:  "asd",
		ChatID:    message.ChatID,
		Operation: network.SYNC_RESPONSE,
	}

	//TODO send sync request to other peer to get the difference between the messages that the other peer knows this peer does not know
	syncRequest := network.Message{
		Id:        uuid.New().String(),
		Timestamp: time.Now().UnixNano(),
		Content:   string(internalMessagesBytes),
		FromUser:  "asd",
		ChatID:    message.ChatID,
		Operation: network.SYNC_REQUEST,
	}

	return nil
}

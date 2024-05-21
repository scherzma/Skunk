package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
	"time"
)

// SyncRequestHandler handles the "SyncRequest" message operation.
type syncRequestHandler struct {
	syncStorage           store.SyncStoragePort
	networkMessageStorage store.NetworkMessageStoragePort
	messageSender         MessageSender
}

// HandleMessage processes the received "SyncRequest" message.
// It retrieves the chat messages from the repository, finds the missing messages
// between the current peer and the other peer, and sends a sync response and a
// sync request to the other peer.
func NewSyncRequestHandler(syncStorage store.SyncStoragePort, networkMessageStorage store.NetworkMessageStoragePort) *syncRequestHandler {
	return &syncRequestHandler{
		syncStorage:           syncStorage,
		networkMessageStorage: networkMessageStorage,
	}
}

func (s *syncRequestHandler) HandleMessage(message network.Message) error {

	chatRepo := p_model.GetNetworkChatsInstance()       // TODO: change
	chatMessageRepo := chatRepo.GetChat(message.ChatID) // TODO: change

	// Structure of the message:
	/*
		{
		  "existingMessageIds": [
			"<message id 1>",
			"<message id 2>",
			...
		  ]
		}
	*/

	var content struct {
		ExistingMessageIDs []string `json:"existingMessageIds"`
	}
	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	// Find difference between "message" already known messages and own messages that the other peer does not know
	missingExternalMessages := chatMessageRepo.GetMissingExternalMessages(content.ExistingMessageIDs)   // TODO: change
	missingInternalMessages := chatMessageRepo.GetMissingInternalMessageIDs(content.ExistingMessageIDs) // TODO: change

	// Convert missingExternalMessages to a JSON string
	externalMessagesBytes, err := json.Marshal(missingExternalMessages)
	if err != nil {
		fmt.Println("Error marshalling missing external messages")
		return err
	}

	// Convert missingInternalMessages to a JSON string
	internalMessagesBytes, err := json.Marshal(missingInternalMessages)
	if err != nil {
		fmt.Println("Error marshalling missing internal messages")
		return err
	}

	// Send the sync response to the other peer
	syncResponse := network.Message{
		Id:              uuid.New().String(),
		Timestamp:       time.Now().UnixNano(),
		Content:         string(externalMessagesBytes),
		SenderID:        chatMessageRepo.GetUsername(), // TODO: change
		ReceiverID:      message.SenderID,              // TODO: change
		SenderAddress:   message.SenderAddress,         // TODO: change
		ReceiverAddress: message.ReceiverAddress,       // TODO: change
		ChatID:          message.ChatID,
		Operation:       network.SYNC_RESPONSE,
	}

	// Send sync request to other peer to get the difference between the messages that the other peer knows this peer does not know
	syncRequest := network.Message{
		Id:              uuid.New().String(),
		Timestamp:       time.Now().UnixNano(),
		Content:         string(internalMessagesBytes),
		SenderID:        chatMessageRepo.GetUsername(),
		ReceiverID:      message.ReceiverID,
		SenderAddress:   message.SenderAddress,
		ReceiverAddress: message.ReceiverAddress,
		ChatID:          message.ChatID,
		Operation:       network.SYNC_REQUEST,
	}

	s.messageSender.SendMessage(syncResponse)
	s.messageSender.SendMessage(syncRequest)

	return nil
}

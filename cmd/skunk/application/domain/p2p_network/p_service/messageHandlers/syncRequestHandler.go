package messageHandlers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
)

// SyncRequestHandler handles the "SyncRequest" message operation.
type syncRequestHandler struct {
	syncStorage           store.SyncStoragePort
	networkMessageStorage store.NetworkMessageStoragePort
	messageSender         MessageSender
}

// NewSyncRequestHandler creates a new instance of syncRequestHandler.
func NewSyncRequestHandler(syncStorage store.SyncStoragePort, networkMessageStorage store.NetworkMessageStoragePort, messageSender MessageSender) *syncRequestHandler {
	return &syncRequestHandler{
		syncStorage:           syncStorage,
		networkMessageStorage: networkMessageStorage,
		messageSender:         messageSender,
	}
}

// HandleMessage processes the received "SyncRequest" message.
func (s *syncRequestHandler) HandleMessage(message network.Message) error {
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

	// Get missing messages
	missingExternalMessageIDs, err := s.syncStorage.GetMissingExternalMessages(message.ChatID, content.ExistingMessageIDs)
	if err != nil {
		fmt.Println("Error getting missing external messages")
		return err
	}

	var missingExternalMessages []network.Message
	missingExternalMessages = make([]network.Message, len(missingExternalMessageIDs))

	for i, messageID := range missingExternalMessageIDs {
		missingExternalMessages[i], err = s.networkMessageStorage.RetrieveMessage(messageID)
		if err != nil {
			return err
		}
	}

	// TODO FIX
	missingInternalMessages, err := s.syncStorage.GetMissingInternalMessages(message.ChatID, content.ExistingMessageIDs)
	if err != nil {
		fmt.Println("Error getting missing internal messages")
		return err
	}

	// Convert missing messages to JSON strings
	externalMessagesBytes, err := json.Marshal(missingExternalMessages)
	if err != nil {
		fmt.Println("Error marshalling missing external messages")
		return err
	}

	internalMessagesBytes, err := json.Marshal(missingInternalMessages)
	if err != nil {
		fmt.Println("Error marshalling missing internal messages")
		return err
	}

	// Create and send the sync response
	syncResponse := network.Message{
		Id:              uuid.New().String(),
		Timestamp:       time.Now().UnixNano(),
		Content:         string(externalMessagesBytes),
		SenderID:        message.ReceiverID,
		ReceiverID:      message.SenderID,
		SenderAddress:   message.ReceiverAddress,
		ReceiverAddress: message.SenderAddress,
		ChatID:          message.ChatID,
		Operation:       network.SYNC_RESPONSE,
	}

	// Create and send the sync request
	syncRequest := network.Message{
		Id:              uuid.New().String(),
		Timestamp:       time.Now().UnixNano(),
		Content:         string(internalMessagesBytes),
		SenderID:        message.ReceiverID,
		ReceiverID:      message.SenderID,
		SenderAddress:   message.ReceiverAddress,
		ReceiverAddress: message.SenderAddress,
		ChatID:          message.ChatID,
		Operation:       network.SYNC_REQUEST,
	}

	s.messageSender.SendMessage(syncResponse)
	s.messageSender.SendMessage(syncRequest)

	return nil
}

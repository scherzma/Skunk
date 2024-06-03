package test

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestSendMessageHandler verifies the sendMessageHandler using a mock handler
func TestSendMessageHandler(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_send_message_handler.db"
	defer os.Remove(dbPath)

	// Initialize storage adapter
	adapter := storageSQLiteAdapter.GetInstance(dbPath)

	// Create a mock chat logic
	mockChatLogic := &MockChatLogic{}

	// Create a sendMessageHandler with the mock chat logic and storage adapter
	sendMessageHandler := messageHandlers.NewSendMessageHandler(mockChatLogic, adapter)

	// Prepare a send message
	messageContent := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	}

	messageContentBytes, err := json.Marshal(messageContent)
	if err != nil {
		t.Fatalf("Error marshalling message content: %v", err)
	}

	sendMessage := network.Message{
		Id:              "msg1",
		Timestamp:       1633029460,
		Content:         string(messageContentBytes),
		SenderID:        "user1",
		ReceiverID:      "user2",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "user2.onion",
		ChatID:          "chat1",
		Operation:       network.SEND_MESSAGE,
	}

	// Handle the send message
	err = sendMessageHandler.HandleMessage(sendMessage)
	if err != nil {
		t.Fatalf("Error handling send message: %v", err)
	}

	// Verify that the message was stored in the database
	retrievedMessage, err := adapter.RetrieveMessage(sendMessage.Id)
	if err != nil {
		t.Fatalf("Error retrieving message: %v", err)
	}
	assert.Equal(t, sendMessage, retrievedMessage)

	// Verify that the message was processed correctly by the chat logic
	assert.Equal(t, "user1", mockChatLogic.LastSenderId)
	assert.Equal(t, "chat1", mockChatLogic.LastChatId)
	assert.Equal(t, "Hello, World!", mockChatLogic.LastMessage)

	fmt.Println("Send message handler test passed")
}

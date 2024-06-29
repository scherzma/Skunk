package test

import (
	"encoding/json"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSendMessageHandler(t *testing.T) {
	dbPath := "test_send_message_handler.db"
	defer os.Remove(dbPath)
	t.Logf("Using temporary database: %s", dbPath)

	adapter := storageSQLiteAdapter.GetInstance(dbPath)
	t.Log("Storage adapter initialized")

	mockChatLogic := &MockChatLogic{}
	t.Log("Mock chat logic created")

	sendMessageHandler := messageHandlers.NewSendMessageHandler(mockChatLogic, adapter)
	t.Log("Send message handler created")

	messageContent := struct {
		Message string `json:"message"`
	}{
		Message: "Hello, World!",
	}

	messageContentBytes, err := json.Marshal(messageContent)
	assert.NoError(t, err, "Error marshalling message content")
	t.Log("Message content prepared")

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
	t.Logf("Send message created: %+v", sendMessage)

	err = sendMessageHandler.HandleMessage(sendMessage)
	assert.NoError(t, err, "Error handling send message")

	retrievedMessage, err := adapter.RetrieveMessage(sendMessage.Id)
	assert.NoError(t, err, "Error retrieving message")
	assert.Equal(t, sendMessage, retrievedMessage, "Retrieved message does not match sent message")
	t.Log("Message successfully stored and retrieved from database")

	assert.Equal(t, "user1", mockChatLogic.LastSenderId, "Unexpected LastSenderId")
	assert.Equal(t, "chat1", mockChatLogic.LastChatId, "Unexpected LastChatId")
	assert.Equal(t, "Hello, World!", mockChatLogic.LastMessage, "Unexpected LastMessage")
	t.Log("Message correctly processed by mock chat logic")

	assert.Contains(t, mockChatLogic.LogEntries, "ReceiveMessage called", "Expected ReceiveMessage to be called")
	t.Log("Send message handler test passed")
}

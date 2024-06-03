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

func TestSetUsernameHandler(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_set_username.db"
	defer os.Remove(dbPath)

	// Initialize storage adapter
	adapter := storageSQLiteAdapter.GetInstance(dbPath)

	// Create a mock chat logic
	mockChatLogic := &MockChatLogic{}

	// Create a setUsernameHandler with the mock chat logic and storage adapter
	setUsernameHandler := messageHandlers.NewSetUsernameHandler(mockChatLogic, adapter)

	// Prepare a set username message
	usernameContent := struct {
		Username string `json:"username"`
	}{
		Username: "CoolUser1",
	}

	usernameContentBytes, err := json.Marshal(usernameContent)
	if err != nil {
		t.Fatalf("Error marshalling username content: %v", err)
	}

	setUsernameMessage := network.Message{
		Id:              "setUsernameMsg1",
		Timestamp:       1633029460,
		Content:         string(usernameContentBytes),
		SenderID:        "user1",
		ReceiverID:      "",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "",
		ChatID:          "chat1",
		Operation:       network.SET_USERNAME,
	}

	// Handle the set username message
	err = setUsernameHandler.HandleMessage(setUsernameMessage)
	if err != nil {
		t.Fatalf("Error handling set username message: %v", err)
	}

	// Verify that the username was stored in the database
	username, err := adapter.GetUsername("user1", "chat1")
	if err != nil {
		t.Fatalf("Error getting username: %v", err)
	}
	assert.Equal(t, "CoolUser1", username)

	// Verify that the username change was processed correctly by the chat logic
	assert.Equal(t, "user1", mockChatLogic.LastSenderId)
	assert.Equal(t, "chat1", mockChatLogic.LastChatId)
	assert.Equal(t, "CoolUser1", mockChatLogic.LastUsername)

	fmt.Println("Set username handler test passed")
}

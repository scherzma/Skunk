package test

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
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

	// Create a mock network connection
	mockNetworkConnection := networkMockAdapter.GetMockConnection()

	// Create a peer and add the mock network connection
	peer := messageHandlers.GetPeerInstance()
	peer.AddNetworkConnection(mockNetworkConnection)

	// Create an invitation for a user, because otherwise the user will not be able to join the chat
	inviteContent := struct {
		ChatID   string                   `json:"chatId"`
		ChatName string                   `json:"chatName"`
		Peers    []store.PublicKeyAddress `json:"peers"`
	}{
		ChatID:   "chat1",
		ChatName: "Cool Chat",
		Peers: []store.PublicKeyAddress{
			{PublicKey: "user3.onion", Address: "user3"},
		},
	}

	inviteContentBytes, err := json.Marshal(inviteContent)
	if err != nil {
		t.Fatalf("Error marshalling invite content: %v", err)
	}

	inviteMessage := network.Message{
		Id:              "inviteMsg1",
		Timestamp:       1633029460,
		Content:         string(inviteContentBytes),
		SenderID:        "user1",
		ReceiverID:      "user2",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "user2.onion",
		ChatID:          "chat1",
		Operation:       network.INVITE_TO_CHAT,
	}

	// Prepare a join chat message
	joinChatMessage := network.Message{
		Id:              "joinMsg1",
		Timestamp:       1633029460,
		Content:         "",
		SenderID:        "user2",
		ReceiverID:      "",
		SenderAddress:   "user2.onion",
		ReceiverAddress: "",
		ChatID:          "chat1",
		Operation:       network.JOIN_CHAT,
	}

	mockNetworkConnection.SendMockNetworkMessageToSubscribers(inviteMessage)
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(joinChatMessage)

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
		SenderID:        "user2",
		ReceiverID:      "",
		SenderAddress:   "user2.onion",
		ReceiverAddress: "",
		ChatID:          "chat1",
		Operation:       network.SET_USERNAME,
	}

	// Create a test chat, so that the peer can set the username
	adapter.CreateChat("chat1", "Test Chat")

	// Handle the set username message
	err = mockNetworkConnection.SendMockNetworkMessageToSubscribers(setUsernameMessage)
	if err != nil {
		t.Fatalf("Error handling set username message: %v", err)
	}

	// Verify that the username was stored in the database
	username, err := adapter.GetUsername("user2", "chat1")
	if err != nil {
		t.Fatalf("Error getting username: %v", err)
	}
	assert.Equal(t, "CoolUser1", username)

	fmt.Println("Set username handler test passed")
}

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

func TestNetworkToPeerFlow(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_network_to_peer_flow.db"
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
			{PublicKey: "user5.onion", Address: "address1"},
			{PublicKey: "user6.onion", Address: "address2"},
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

	mockNetworkConnection.SendMockNetworkMessageToSubscribers(inviteMessage)
	// End of the invite flow

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

	// Send the join chat message to trigger the join process
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(joinChatMessage)

	// Verify that the peer was added to the chat in the database
	usersInChat, err := adapter.GetUsersInChat("chat1")
	if err != nil {
		t.Fatalf("Error getting users in chat: %v", err)
	}

	found := false
	for _, user := range usersInChat {
		if user.UserId == "user2" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Expected to find user 'user1' in chat 'chat1', but did not")
	}

	fmt.Println("Network to peer flow test passed")
}

// TestJoinChatHandler verifies the joinChatHandler using a mock handler
func TestJoinChatHandler(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_join_chat_handler.db"
	defer os.Remove(dbPath)

	// Initialize storage adapter
	adapter := storageSQLiteAdapter.GetInstance(dbPath)

	// Create a mock chat logic
	mockChatLogic := &MockChatLogic{}

	// Create a joinChatHandler with the mock chat logic and storage adapter
	joinChatHandler := messageHandlers.NewJoinChatHandler(mockChatLogic, adapter)

	// Prepare a join chat message
	joinChatMessage := network.Message{
		Id:              "joinMsg1",
		Timestamp:       1633029460,
		Content:         "",
		SenderID:        "user1",
		ReceiverID:      "",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "",
		ChatID:          "chat1",
		Operation:       network.JOIN_CHAT,
	}

	// Handle the join chat message
	err := joinChatHandler.HandleMessage(joinChatMessage)
	if err != nil {
		t.Fatalf("Error handling join chat message: %v", err)
	}

	// Verify that the peer joined the chat
	assert.Equal(t, "user1", mockChatLogic.LastSenderId)
	assert.Equal(t, "chat1", mockChatLogic.LastChatId)

	fmt.Println("Join chat handler test passed")
}

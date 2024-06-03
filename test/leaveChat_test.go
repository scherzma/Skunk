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

// TestLeaveChatHandler verifies the leaveChatHandler using a mock handler
func TestLeaveChatHandler(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_leave_chat_handler.db"
	defer os.Remove(dbPath)

	// Initialize storage adapter
	adapter := storageSQLiteAdapter.GetInstance(dbPath)

	// Create a mock chat logic
	mockChatLogic := &MockChatLogic{}

	// Create a leaveChatHandler with the mock chat logic and storage adapter
	leaveChatHandler := messageHandlers.NewLeaveChatHandler(mockChatLogic, adapter)

	// Prepare a leave chat message
	leaveChatMessage := network.Message{
		Id:              "leaveMsg1",
		Timestamp:       1633029460,
		Content:         "",
		SenderID:        "user1",
		ReceiverID:      "",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "",
		ChatID:          "chat1",
		Operation:       network.LEAVE_CHAT,
	}

	// Handle the leave chat message
	err := leaveChatHandler.HandleMessage(leaveChatMessage)
	if err != nil {
		t.Fatalf("Error handling leave chat message: %v", err)
	}

	// Verify that the peer left the chat
	assert.Equal(t, "user1", mockChatLogic.LastSenderId)
	assert.Equal(t, "chat1", mockChatLogic.LastChatId)

	fmt.Println("Leave chat handler test passed")
}

// TestPeerLeavesChat tests the complete flow from network to peer
func TestPeerLeavesChat(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_peer_leaves_chat.db"
	defer os.Remove(dbPath)

	// Initialize storage adapter
	adapter := storageSQLiteAdapter.GetInstance(dbPath)

	// Create a mock network connection
	mockNetworkConnection := networkMockAdapter.GetMockConnection()

	// Create a peer and add the mock network connection
	peer := messageHandlers.GetPeerInstance()
	peer.AddNetworkConnection(mockNetworkConnection)

	// Prepare a leave chat message
	leaveChatMessage := network.Message{
		Id:              "leaveMsg1",
		Timestamp:       1633029460,
		Content:         "",
		SenderID:        "user2",
		ReceiverID:      "",
		SenderAddress:   "user2.onion",
		ReceiverAddress: "",
		ChatID:          "chat1",
		Operation:       network.LEAVE_CHAT,
	}

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
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(leaveChatMessage)

	// Verify that the peer was removed from the chat in the database
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

	if found {
		t.Fatalf("Expected not to find user 'user2' in chat 'chat1', but did")
	}

	fmt.Println("Peer leaves chat flow test passed")
}

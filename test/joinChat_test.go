package test

import (
	"encoding/json"
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
	t.Logf("Using temporary database: %s", dbPath)

	// Initialize storage adapter
	adapter := storageSQLiteAdapter.GetInstance(dbPath)
	t.Log("Storage adapter initialized")

	// Create a mock network connection
	mockNetworkConnection := networkMockAdapter.GetMockConnection()
	t.Log("Mock network connection created")

	// Create a peer and add the mock network connection
	peer := messageHandlers.GetPeerInstance()
	peer.AddNetworkConnection(mockNetworkConnection)
	t.Log("Peer instance created and mock network connection added")

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
	assert.NoError(t, err, "Error marshalling invite content")
	t.Log("Invite content prepared")

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
	t.Logf("Invite message created: %+v", inviteMessage)

	mockNetworkConnection.SendMockNetworkMessageToSubscribers(inviteMessage)
	t.Log("Invite message sent to subscribers")

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
	t.Logf("Join chat message created: %+v", joinChatMessage)

	// Send the join chat message to trigger the join process
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(joinChatMessage)
	t.Log("Join chat message sent to subscribers")

	// Verify that the peer was added to the chat in the database
	usersInChat, err := adapter.GetUsersInChat("chat1")
	assert.NoError(t, err, "Error getting users in chat")
	t.Logf("Users in chat: %+v", usersInChat)

	found := false
	for _, user := range usersInChat {
		if user.UserId == "user2" {
			found = true
			break
		}
	}

	assert.True(t, found, "Expected to find user 'user2' in chat 'chat1', but did not")

	t.Log("Network to peer flow test passed")
}

// TestJoinChatHandler verifies the joinChatHandler using a mock handler
func TestJoinChatHandler(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_join_chat_handler.db"
	defer os.Remove(dbPath)
	t.Logf("Using temporary database: %s", dbPath)

	// Initialize storage adapter
	adapter := storageSQLiteAdapter.GetInstance(dbPath)
	t.Log("Storage adapter initialized")

	// Create a mock chat logic
	mockChatLogic := &MockChatLogic{}
	t.Log("Mock chat logic created")

	// Create a joinChatHandler with the mock chat logic and storage adapter
	joinChatHandler := messageHandlers.NewJoinChatHandler(mockChatLogic, adapter)
	t.Log("Join chat handler created")

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
	t.Logf("Join chat message created: %+v", joinChatMessage)

	// Handle the join chat message
	err := joinChatHandler.HandleMessage(joinChatMessage)
	assert.NoError(t, err, "Error handling join chat message")

	// Verify that the peer joined the chat
	assert.Equal(t, "user1", mockChatLogic.LastSenderId, "Unexpected LastSenderId")
	assert.Equal(t, "chat1", mockChatLogic.LastChatId, "Unexpected LastChatId")

	t.Log("Join chat handler test passed")
}

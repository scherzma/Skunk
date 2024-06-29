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

// TestLeaveChatHandler verifies the leaveChatHandler using a mock handler
func TestLeaveChatHandler(t *testing.T) {
	dbPath := "test_leave_chat_handler.db"
	defer os.Remove(dbPath)
	t.Logf("Using temporary database: %s", dbPath)

	adapter := storageSQLiteAdapter.GetInstance(dbPath)
	t.Log("Storage adapter initialized")

	mockChatLogic := &MockChatLogic{}
	t.Log("Mock chat logic created")

	leaveChatHandler := messageHandlers.NewLeaveChatHandler(mockChatLogic, adapter)
	t.Log("Leave chat handler created")

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
	t.Logf("Leave chat message created: %+v", leaveChatMessage)

	err := leaveChatHandler.HandleMessage(leaveChatMessage)
	assert.NoError(t, err, "Error handling leave chat message")

	assert.Equal(t, "user1", mockChatLogic.LastSenderId, "Unexpected LastSenderId")
	assert.Equal(t, "chat1", mockChatLogic.LastChatId, "Unexpected LastChatId")

	t.Log("Leave chat handler test passed")
}

// TestPeerLeavesChat tests the complete flow from network to peer
func TestPeerLeavesChat(t *testing.T) {
	dbPath := "test_peer_leaves_chat.db"
	defer os.Remove(dbPath)
	t.Logf("Using temporary database: %s", dbPath)

	adapter := storageSQLiteAdapter.GetInstance(dbPath)
	t.Log("Storage adapter initialized")

	mockNetworkConnection := networkMockAdapter.GetMockConnection()
	t.Log("Mock network connection created")

	peer := messageHandlers.GetPeerInstance()
	peer.AddNetworkConnection(mockNetworkConnection)
	t.Log("Peer instance created and mock network connection added")

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
	t.Logf("Leave chat message created: %+v", leaveChatMessage)

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

	mockNetworkConnection.SendMockNetworkMessageToSubscribers(inviteMessage)
	t.Log("Invite message sent")
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(joinChatMessage)
	t.Log("Join chat message sent")
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(leaveChatMessage)
	t.Log("Leave chat message sent")

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

	assert.False(t, found, "Expected not to find user 'user2' in chat 'chat1', but did")

	t.Log("Peer leaves chat flow test passed")
}

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

func TestSetUsernameHandler(t *testing.T) {
	dbPath := "test_set_username.db"
	defer os.Remove(dbPath)
	t.Logf("Using temporary database: %s", dbPath)

	adapter := storageSQLiteAdapter.GetInstance(dbPath)
	t.Log("Storage adapter initialized")

	mockNetworkConnection := networkMockAdapter.GetMockConnection()
	t.Log("Mock network connection created")

	peer := messageHandlers.GetPeerInstance()
	peer.AddNetworkConnection(mockNetworkConnection)
	t.Log("Peer instance created and mock network connection added")

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

	usernameContent := struct {
		Username string `json:"username"`
	}{
		Username: "CoolUser1",
	}

	usernameContentBytes, err := json.Marshal(usernameContent)
	assert.NoError(t, err, "Error marshalling username content")
	t.Log("Username content prepared")

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
	t.Logf("Set username message created: %+v", setUsernameMessage)

	err = adapter.CreateChat("chat1", "Test Chat")
	assert.NoError(t, err, "Error creating test chat")
	t.Log("Test chat created")

	err = mockNetworkConnection.SendMockNetworkMessageToSubscribers(setUsernameMessage)
	assert.NoError(t, err, "Error handling set username message")
	t.Log("Set username message sent")

	username, err := adapter.GetUsername("user2", "chat1")
	assert.NoError(t, err, "Error getting username")
	assert.Equal(t, "CoolUser1", username, "Unexpected username")
	t.Log("Username successfully set and retrieved from database")

	t.Log("Set username handler test passed")
}

package test

import (
	"encoding/json"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
	"os"
	"testing"
)

func TestInviteToChatHandler(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_invite_to_chat.db"
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

	// Prepare an invite to chat message
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
	t.Log("Mock network message sent")

	// Verify that the invitation was stored in the database
	invitations, err := adapter.GetInvitations("user2")
	if err != nil {
		t.Fatalf("Error getting invitations: %v", err)
	}
	t.Logf("Retrieved invitations for user2: %v", invitations)

	found := false
	for _, invitation := range invitations {
		if invitation == "chat1" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Expected to find invitation to 'chat1' for 'user2', but did not")
	}

	t.Log("Invite to chat test passed")
}

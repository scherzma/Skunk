package test

import (
	"encoding/json"
	"fmt"
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

	// Initialize storage adapter
	adapter := storageSQLiteAdapter.GetInstance(dbPath)

	// Create a mock network connection
	mockNetworkConnection := networkMockAdapter.GetMockConnection()

	// Create a peer and add the mock network connection
	peer := messageHandlers.GetPeerInstance()
	peer.AddNetworkConnection(mockNetworkConnection)

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

	//inviteBytes := string(inviteContentBytes)
	//fmt.Println(inviteBytes)

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

	// Verify that the invitation was stored in the database
	invitations, err := adapter.GetInvitations("user2")
	if err != nil {
		t.Fatalf("Error getting invitations: %v", err)
	}

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

	fmt.Println("Invite to chat test passed")
}

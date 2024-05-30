package test

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"os"
	"testing"
)

func TestSyncRequestHandler(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_sync_request.db"
	defer os.Remove(dbPath)

	// Initialize storage adapter
	adapter := storageSQLiteAdapter.GetInstance(dbPath)

	// Create a mock network connection
	mockNetworkConnection := networkMockAdapter.GetMockConnection()

	// Create a peer and add the mock network connection
	peer := messageHandlers.GetPeerInstance()
	peer.AddNetworkConnection(mockNetworkConnection)

	// Create some internal messages
	internalMessages := []network.Message{
		{
			Id:              "msgMsg3",
			Timestamp:       1633029448,
			Content:         "LOOOOOOOOOOOOOOOOOOOOOOOL",
			SenderID:        "user3",
			ReceiverID:      "user4",
			SenderAddress:   "user3.onion",
			ReceiverAddress: "user4.onion",
			ChatID:          "chat1",
			Operation:       network.SEND_MESSAGE,
		},
		{
			Id:              "msgMsg4",
			Timestamp:       1633029448,
			Content:         "WOOW",
			SenderID:        "user3",
			ReceiverID:      "user4",
			SenderAddress:   "user3.onion",
			ReceiverAddress: "user4.onion",
			ChatID:          "chat1",
			Operation:       network.SEND_MESSAGE,
		},
		{
			Id:              "msgMsg5",
			Timestamp:       1633029448,
			Content:         "WOLOLOW",
			SenderID:        "user3",
			ReceiverID:      "user4",
			SenderAddress:   "user3.onion",
			ReceiverAddress: "user4.onion",
			ChatID:          "chat2",
			Operation:       network.SEND_MESSAGE,
		},
	}

	// Store internal messages
	for _, msg := range internalMessages {
		err := adapter.StoreMessage(msg)
		if err != nil {
			t.Errorf("Error storing internal message: %v", err)
		}
	}

	// Prepare a sync request message
	testSyncMessage := network.Message{
		Id:              "syncMsg2",
		Timestamp:       1633029446,
		Content:         "{\"existingMessageIds\": [\"msgMsg5\",\"msg2\"]}",
		SenderID:        "user1",
		ReceiverID:      "user2",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "user2.onion",
		ChatID:          "chat1",
		Operation:       network.SYNC_REQUEST,
	}

	adapter.PeerJoinedChat(3214523465, "user1", "chat1")

	// Send the sync request message to trigger the sync process
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(testSyncMessage)

	fmt.Println(mockNetworkConnection.LastSent)

}

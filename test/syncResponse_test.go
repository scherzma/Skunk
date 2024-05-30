package test

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSyncResponseHandler(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test_sync_response.db"
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

	// Create some messages to be synced
	messagesToSync := []network.Message{
		{
			Id:              "msg1",
			Timestamp:       1633029445,
			Content:         "Hello World!",
			SenderID:        "user1",
			ReceiverID:      "user2",
			SenderAddress:   "user1.onion",
			ReceiverAddress: "user2.onion",
			ChatID:          "chat1",
			Operation:       network.SEND_MESSAGE,
		},
		{
			Id:              "msg2",
			Timestamp:       1633029446,
			Content:         "Hello Again!",
			SenderID:        "user2",
			ReceiverID:      "user1",
			SenderAddress:   "user2.onion",
			ReceiverAddress: "user1.onion",
			ChatID:          "chat1",
			Operation:       network.SEND_MESSAGE,
		},
	}

	// Convert missing messages to JSON strings
	syncResponseContent, err := json.Marshal(messagesToSync)
	if err != nil {
		fmt.Println("Error marshalling missing external messages")
	}

	// Prepare a sync response message
	testSyncResponseMessage := network.Message{
		Id:              "syncMsg1",
		Timestamp:       1633029450,
		Content:         string(syncResponseContent),
		SenderID:        "user2",
		ReceiverID:      "user1",
		SenderAddress:   "user2.onion",
		ReceiverAddress: "user1.onion",
		ChatID:          "chat1",
		Operation:       network.SYNC_RESPONSE,
	}

	adapter.PeerJoinedChat(3214523465, "user2", "chat1")

	// Send the sync response message to trigger the sync process
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(testSyncResponseMessage)

	// Retrieve messages from the database and verify they were correctly stored
	for _, originalMessage := range messagesToSync {
		retrievedMessage, err := adapter.RetrieveMessage(originalMessage.Id)
		if err != nil {
			t.Errorf("Error retrieving message: %v", err)
		}

		assert.Equal(t, originalMessage, retrievedMessage, "Stored message does not match original message")
	}

	fmt.Println("Sync response test passed")
}

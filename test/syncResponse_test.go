package test

import (
	"encoding/json"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSyncResponseHandler(t *testing.T) {
	dbPath := "test_sync_response.db"
	defer os.Remove(dbPath)
	t.Logf("Using temporary database: %s", dbPath)

	adapter := storageSQLiteAdapter.GetInstance(dbPath)
	t.Log("Storage adapter initialized")

	mockNetworkConnection := networkMockAdapter.GetMockConnection()
	t.Log("Mock network connection created")

	peer := messageHandlers.GetPeerInstance()
	err := peer.AddNetworkConnection(mockNetworkConnection)
	assert.NoError(t, err, "Error adding mock network connection to peer")
	t.Log("Peer instance created and mock network connection added")

	t.Run("StoreInternalMessages", func(t *testing.T) {
		internalMessages := getInternalMessages_Response()
		for _, msg := range internalMessages {
			err := adapter.StoreMessage(msg)
			assert.NoError(t, err, "Error storing internal message")
		}
		t.Log("Internal messages stored successfully")
	})

	t.Run("PrepareSyncResponseMessage", func(t *testing.T) {
		messagesToSync := getMessagesToSync()
		syncResponseContent, err := json.Marshal(messagesToSync)
		assert.NoError(t, err, "Error marshalling sync response content")

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
		t.Logf("Sync response message prepared: %+v", testSyncResponseMessage)

		err = adapter.PeerJoinedChat(3214523465, "user2", "chat1")
		assert.NoError(t, err, "Error adding peer to chat")

		err = mockNetworkConnection.SendMockNetworkMessageToSubscribers(testSyncResponseMessage)
		assert.NoError(t, err, "Error sending sync response message")
		t.Log("Sync response message sent successfully")
	})

	t.Run("VerifySyncedMessages", func(t *testing.T) {
		messagesToSync := getMessagesToSync()
		for _, originalMessage := range messagesToSync {
			retrievedMessage, err := adapter.RetrieveMessage(originalMessage.Id)
			assert.NoError(t, err, "Error retrieving message")
			assert.Equal(t, originalMessage, retrievedMessage, "Stored message does not match original message")
		}
		t.Log("All synced messages verified successfully")
	})

	t.Log("Sync response test passed")
}

func getInternalMessages_Response() []network.Message {
	return []network.Message{
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
}

func getMessagesToSync() []network.Message {
	return []network.Message{
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
}

package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestSyncRequestHandler(t *testing.T) {
	dbPath := "test_sync_request.db"
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
		internalMessages := getInternalMessages()
		for _, msg := range internalMessages {
			err := adapter.StoreMessage(msg)
			assert.NoError(t, err, "Error storing internal message")
		}
		t.Log("Internal messages stored successfully")
	})

	t.Run("SendSyncRequestMessage", func(t *testing.T) {
		testSyncMessage := network.Message{
			Id:              "syncMsg2",
			Timestamp:       1633029446,
			Content:         "{\"existingMessageIds\": [\"msgMsg9\",\"msg2\",\"msgMsg3\"]}",
			SenderID:        "user1",
			ReceiverID:      "user2",
			SenderAddress:   "user1.onion",
			ReceiverAddress: "user2.onion",
			ChatID:          "chat1",
			Operation:       network.SYNC_REQUEST,
		}
		t.Logf("Sync request message created: %+v", testSyncMessage)

		err := adapter.PeerJoinedChat(3214523465, "user1", "chat1")
		assert.NoError(t, err, "Error adding peer to chat")
		t.Log("Peer joined chat successfully")

		err = mockNetworkConnection.SendMockNetworkMessageToSubscribers(testSyncMessage)
		assert.NoError(t, err, "Error sending sync request message")
		t.Log("Sync request message sent successfully")
	})

	t.Run("VerifySyncResponse", func(t *testing.T) {
		t.Logf("Last sent message: %+v", mockNetworkConnection.LastSent)
		assert.Equal(t, "user2", mockNetworkConnection.LastSent.SenderID, "Unexpected SenderID")
		assert.Equal(t, "chat1", mockNetworkConnection.LastSent.ChatID, "Unexpected ChatID")
		assert.Equal(t, "user1", mockNetworkConnection.LastSent.ReceiverID, "Unexpected ReceiverID")
		assert.Equal(t, "user1.onion", mockNetworkConnection.LastSent.ReceiverAddress, "Unexpected ReceiverAddress")
		assert.Equal(t, "user2.onion", mockNetworkConnection.LastSent.SenderAddress, "Unexpected SenderAddress")
		assert.Equal(t, network.SYNC_REQUEST, mockNetworkConnection.LastSent.Operation, "Unexpected Operation")
		t.Log("Sync response verified successfully")
	})
}

func getInternalMessages() []network.Message {
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

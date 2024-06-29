package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestStorageSQLiteAdapter(t *testing.T) {
	dbPath := "test.db"
	defer os.Remove(dbPath)
	t.Logf("Using temporary database: %s", dbPath)

	adapter := storageSQLiteAdapter.GetInstance(dbPath)
	t.Log("Storage adapter initialized")

	testMessages := getTestMessages()

	t.Run("StoreAndRetrieveMessages", func(t *testing.T) {
		for _, msg := range testMessages {
			err := adapter.StoreMessage(msg)
			assert.NoError(t, err, "Error storing message")
		}
		t.Log("Messages stored")

		for _, msg := range testMessages {
			retrieved, err := adapter.RetrieveMessage(msg.Id)
			assert.NoError(t, err, "Error retrieving message")
			assert.Equal(t, msg, retrieved, "Retrieved message does not match stored message")
		}
		t.Log("All messages successfully retrieved and matched")
	})

	t.Run("GetChatMessages", func(t *testing.T) {
		chatMessages, err := adapter.GetChatMessages("chat1")
		assert.NoError(t, err, "Error getting chat messages")
		assert.Equal(t, len(testMessages), len(chatMessages), "Unexpected number of chat messages")
		t.Log("GetChatMessages successful")
	})

	t.Run("SetPeerUsername", func(t *testing.T) {
		err := adapter.SetPeerUsername("CoolUser1", "user1", "chat1")
		assert.NoError(t, err, "Error setting peer username")
		t.Log("SetPeerUsername successful")
	})

	t.Run("PeerJoinedChat", func(t *testing.T) {
		err := adapter.PeerJoinedChat(340982203948, "user1", "chat1")
		assert.NoError(t, err, "Error adding peer to chat")
		t.Log("PeerJoinedChat successful")
	})

	t.Run("PeerLeftChat", func(t *testing.T) {
		err := adapter.PeerLeftChat("user1", "chat1")
		assert.NoError(t, err, "Error removing peer from chat")
		t.Log("PeerLeftChat successful")
	})

	t.Run("CreateAndGetChats", func(t *testing.T) {
		err := adapter.CreateChat("chat2", "Test Chat")
		assert.NoError(t, err, "Error creating chat")

		chats, err := adapter.GetChats()
		assert.NoError(t, err, "Error getting chats")

		foundChat := false
		for _, chat := range chats {
			if chat.ChatId == "chat2" && chat.ChatName == "Test Chat" {
				foundChat = true
				break
			}
		}
		assert.True(t, foundChat, "Created chat not found in GetChats")
		t.Log("CreateChat and GetChats successful")
	})

	t.Run("InvitedToChat", func(t *testing.T) {
		err := adapter.InvitedToChat("msg1", []store.PublicKeyAddress{
			{PublicKey: "user5.onion", Address: "address1"},
			{PublicKey: "user6.onion", Address: "address2"},
		})
		assert.NoError(t, err, "Error inviting to chat")
		t.Log("InvitedToChat successful")
	})
}

func getTestMessages() []network.Message {
	return []network.Message{
		{
			Id:              "msg1",
			Timestamp:       1620000000,
			Content:         "{\"message\": \"Hey everyone!\"}",
			SenderID:        "user1",
			ReceiverID:      "user2",
			SenderAddress:   "user1.onion",
			ReceiverAddress: "user2.onion",
			ChatID:          "chat1",
			Operation:       network.SEND_MESSAGE,
		},
		{
			Id:              "sync1",
			Timestamp:       1620000060,
			Content:         "{\"existingMessageIds\": [\"msg1\"]}",
			SenderID:        "user2",
			ReceiverID:      "user1",
			SenderAddress:   "user2.onion",
			ReceiverAddress: "user1.onion",
			ChatID:          "chat1",
			Operation:       network.SYNC_REQUEST,
		},
		{
			Id:              "sync2",
			Timestamp:       1620000120,
			Content:         "[{\"Id\":\"msg2\",\"Timestamp\":1620000090,\"Content\":\"{\\\"message\\\": \\\"Hi User1!\\\"}\",\"SenderID\":\"user3\",\"ReceiverID\":\"user1\",\"SenderAddress\":\"user3.onion\",\"ReceiverAddress\":\"user1.onion\",\"ChatID\":\"chat1\",\"Operation\":0}]",
			SenderID:        "user1",
			ReceiverID:      "user2",
			SenderAddress:   "user1.onion",
			ReceiverAddress: "user2.onion",
			ChatID:          "chat1",
			Operation:       network.SYNC_RESPONSE,
		},
		{
			Id:              "join1",
			Timestamp:       1620000180,
			Content:         "",
			SenderID:        "user4",
			ReceiverID:      "",
			SenderAddress:   "user4.onion",
			ReceiverAddress: "",
			ChatID:          "chat1",
			Operation:       network.JOIN_CHAT,
		},
		{
			Id:              "leave1",
			Timestamp:       1620000240,
			Content:         "",
			SenderID:        "user3",
			ReceiverID:      "",
			SenderAddress:   "user3.onion",
			ReceiverAddress: "",
			ChatID:          "chat1",
			Operation:       network.LEAVE_CHAT,
		},
		{
			Id:              "invite1",
			Timestamp:       1620000300,
			Content:         "{\"chatId\": \"chat1\", \"chatName\": \"Cool Chat\", \"peers\": [\"user5.onion\", \"user6.onion\"]}",
			SenderID:        "user1",
			ReceiverID:      "",
			SenderAddress:   "user1.onion",
			ReceiverAddress: "",
			ChatID:          "chat1",
			Operation:       network.INVITE_TO_CHAT,
		},
		{
			Id:              "file1",
			Timestamp:       1620000360,
			Content:         "{\"fileContent\": \"aGVsbG8gd29ybGQ=\"}",
			SenderID:        "user2",
			ReceiverID:      "",
			SenderAddress:   "user2.onion",
			ReceiverAddress: "",
			ChatID:          "chat1",
			Operation:       network.SEND_FILE,
		},
		{
			Id:              "setuser1",
			Timestamp:       1620000420,
			Content:         "{\"username\": \"CoolUser1\"}",
			SenderID:        "user1",
			ReceiverID:      "",
			SenderAddress:   "user1.onion",
			ReceiverAddress: "",
			ChatID:          "chat1",
			Operation:       network.SET_USERNAME,
		},
		{
			Id:              "test1",
			Timestamp:       1620000480,
			Content:         "This is a test message",
			SenderID:        "user1",
			ReceiverID:      "user2",
			SenderAddress:   "user1.onion",
			ReceiverAddress: "user2.onion",
			ChatID:          "chat1",
			Operation:       network.TEST_MESSAGE,
		},
		{
			Id:              "test2",
			Timestamp:       1620000540,
			Content:         "This is another test message",
			SenderID:        "user2",
			ReceiverID:      "user4",
			SenderAddress:   "user2.onion",
			ReceiverAddress: "user4.onion",
			ChatID:          "chat1",
			Operation:       network.TEST_MESSAGE_2,
		},
	}
}

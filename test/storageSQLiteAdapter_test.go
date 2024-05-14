package test

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"os"
	"reflect"
	"testing"
)

func TestStorageSQLiteAdapter(t *testing.T) {
	// Create a temporary database for testing
	dbPath := "test.db"
	defer os.Remove(dbPath)

	adapter := storageSQLiteAdapter.NewStorageSQLiteAdapter(dbPath)

	testMessages := []network.Message{
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

	// Store the test messages
	for _, msg := range testMessages {
		err := adapter.StoreMessage(msg)
		if err != nil {
			t.Errorf("Error storing message: %v", err)
		}
	}

	fmt.Println("Messages stored")
	// Retrieve the messages and compare
	for _, msg := range testMessages {
		retrieved, err := adapter.RetrieveMessage(msg.Id)
		if err != nil {
			t.Errorf("Error retrieving message: %v", err)
		}
		fmt.Println(retrieved)
		fmt.Println(msg)
		if !reflect.DeepEqual(msg, retrieved) {
			t.Errorf("Retrieved message does not match stored message")
		}
	}

	// Test GetChatMessages
	chatMessages, err := adapter.GetChatMessages("chat1")
	if err != nil {
		t.Errorf("Error getting chat messages: %v", err)
	}

	if len(chatMessages) != len(testMessages) {
		t.Errorf("Expected %d chat messages, got %d", len(testMessages), len(chatMessages))
	}

	// Test SetPeerUsername
	err = adapter.SetPeerUsername("CoolUser1", "user1", "chat1")
	if err != nil {
		t.Errorf("Error setting peer username: %v", err)
	}

	username, err := adapter.GetUsername("user1", "chat1")
	if err != nil {
		t.Errorf("Error getting username: %v", err)
	}

	if username != "CoolUser1" {
		t.Errorf("Expected username 'CoolUser1', got '%s'", username)
	}

	// Additional tests can be added for other methods...
}

package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"testing"
)

func TestSyncRequestHandler(t *testing.T) {
	// Create a mock network connection
	testMessage := network.Message{
		Id:              "8888",
		Timestamp:       1633029445,
		Content:         "Hello World!",
		SenderID:        "user1",
		ReceiverID:      "user2",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "user2.onion",
		ChatID:          "chat1",
		Operation:       network.TEST_MESSAGE,
	}

	peer := messageHandlers.GetPeerInstance()

	mockNetworkConnection := networkMockAdapter.GetMockConnection()
	peer.AddNetworkConnection(mockNetworkConnection)

	mockNetworkConnection.SendMockNetworkMessageToSubscribers(testMessage)

	testSyncMessage := network.Message{
		Id:              "12345",
		Timestamp:       1633029446,
		Content:         "{\"existingMessageIds\": [\"<message id 1>\",\"<message id 2>\"]}",
		SenderID:        "user1",
		ReceiverID:      "user2",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "user2.onion",
		ChatID:          "chat1",
		Operation:       network.SYNC_REQUEST,
	}

	internalMessage := network.Message{
		Id:              "internalMessage123!",
		Timestamp:       1633029448,
		Content:         "LOOOOOOOOOOOOOOOOOOOOOOOL",
		SenderID:        "user3",
		ReceiverID:      "user4",
		SenderAddress:   "user3.onion",
		ReceiverAddress: "user4.onion",
		ChatID:          "chat1",
		Operation:       network.SYNC_REQUEST,
	}

	internalMessage2 := network.Message{
		Id:              "internalMessage2",
		Timestamp:       1633029448,
		Content:         "WOOW",
		SenderID:        "user3",
		ReceiverID:      "user4",
		SenderAddress:   "user3.onion",
		ReceiverAddress: "user4.onion",
		ChatID:          "chat1",
		Operation:       network.SYNC_REQUEST,
	}

	internalMessage3 := network.Message{
		Id:              "internalMessage3",
		Timestamp:       1633029448,
		Content:         "WOLOLOW",
		SenderID:        "user3",
		ReceiverID:      "user4",
		SenderAddress:   "user3.onion",
		ReceiverAddress: "user4.onion",
		ChatID:          "chat2",
		Operation:       network.SYNC_REQUEST,
	}

	chat := p_model.GetNetworkChatsInstance().GetChat(internalMessage.ChatID)
	chat.AddMessage(internalMessage)
	chat.AddMessage(internalMessage2)

	chat2 := p_model.GetNetworkChatsInstance().GetChat(internalMessage3.ChatID)
	chat2.AddMessage(internalMessage3)

	mockNetworkConnection.SendMockNetworkMessageToSubscribers(testSyncMessage)

}

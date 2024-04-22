package main

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

func main() {
	// Create a new SQLite storage for the message queue
	/*
		storeMessageQueueSQLite := storage.StoreMessageQueueSQLite{}
		storeMessageQueueSQLite.StoreMessageQueue("aasdf")

		user1 := c_model.User{Username: "user1", UserId: "1"}
		user2 := c_model.User{Username: "user2", UserId: "2"}
		user3 := c_model.User{Username: "user3", UserId: "3"}
		user4 := c_model.User{Username: "user4", UserId: "4"}

		message1 := c_model.Message{
			Id:        "1",
			Timestamp: 1633029442,
			Content:   "Hello, user2!",
			FromUser:      user1,
			chatID:        user2,
			Operation: c_model.SEND_MESSAGE,
		}
		message2 := c_model.Message{
			Id:        "2",
			Timestamp: 1633029443,
			Content:   "Hello, user1!",
			FromUser:      user2,
			chatID:        user1,
			Operation: c_model.SEND_MESSAGE,
		}
		message3 := c_model.Message{
			Id:        "3",
			Timestamp: 1633029444,
			Content:   "Hello, user3!",
			FromUser:      user1,
			chatID:        user3,
			Operation: c_model.SEND_MESSAGE,
		}
		message4 := c_model.Message{
			Id:        "4",
			Timestamp: 1633029445,
			Content:   "Hello, user4!",
			FromUser:      user2,
			chatID:        user4,
			Operation: c_model.SEND_MESSAGE,
		}

		chats := p_model.NewNetworkChatMessages()
		chats.AddMessages([]c_model.Message{message1, message2, message3, message4})

		fmt.Printf("Chats: %v\n", chats.GetMessages())
	*/

	testMessage := network.Message{
		Id:        "8888",
		Timestamp: 1633029445,
		Content:   "Hello asdfasdfWorld!",
		FromUser:  "asd",
		ChatID:    "asdf",
		Operation: network.TEST_MESSAGE,
	}

	peer := messageHandlers.GetPeerInstance()

	mockNetworkConnection := networkMockAdapter.GetMockConnection()
	peer.AddNetworkConnection(mockNetworkConnection)

	mockNetworkConnection.SendMockNetworkMessageToSubscribers(testMessage)

	testSyncMessage := network.Message{
		Id:        "12345",
		Timestamp: 1633029446,
		Content:   "{\"existingMessageIds\": [\"<message id 1>\",\"<message id 2>\"]}",
		FromUser:  "asd",
		ChatID:    "asdf",
		Operation: network.SYNC_REQUEST,
	}

	internalMessage := network.Message{
		Id:        "internalMessage123!",
		Timestamp: 1633029448,
		Content:   "LOOOOOOOOOOOOOOOOOOOOOOOL",
		FromUser:  "as23d",
		ChatID:    "asdf",
		Operation: network.SYNC_REQUEST,
	}

	internalMessage2 := network.Message{
		Id:        "internalMessage2",
		Timestamp: 1633029448,
		Content:   "WOOW",
		FromUser:  "as23d",
		ChatID:    "asdf",
		Operation: network.SYNC_REQUEST,
	}

	internalMessage3 := network.Message{
		Id:        "internalMessage3",
		Timestamp: 1633029448,
		Content:   "WOLOLOW",
		FromUser:  "as23d",
		ChatID:    "asdf1",
		Operation: network.SYNC_REQUEST,
	}

	chat := p_model.GetNetworkChatsInstance().GetChat(internalMessage.ChatID)
	chat.AddMessage(internalMessage)
	chat.AddMessage(internalMessage2)

	chat2 := p_model.GetNetworkChatsInstance().GetChat(internalMessage3.ChatID)
	chat2.AddMessage(internalMessage3)

	mockNetworkConnection.SendMockNetworkMessageToSubscribers(testSyncMessage)

	//{"Id":"internalMessage123!","Timestamp":1633029448,"Content":"LOOOOOOOOOOOOOOOOOOOOOOOL","FromUser":"as23d","ChatID":"asdf","Operation":1},{"Id":"internalMessage2","Timestamp":1633029448,"Content":"WOOW","FromUser":"as23d","ChatID":"asdf","Operation":1}
	mockSyncResponse := network.Message{
		Id:        "mockSyncResponse",
		Timestamp: 1633029449,
		Content:   "[{\"Id\":\"internalMessage123!\",\"Timestamp\":1633029448,\"Content\":\"LOOOOOOOOOOOOOOOOOOOOOOOL\",\"FromUser\":\"as23d\",\"ChatID\":\"asdf\",\"Operation\":1},{\"Id\":\"internalMessage2\",\"Timestamp\":1633029448,\"Content\":\"WOOW\",\"FromUser\":\"as23d\",\"ChatID\":\"asdf\",\"Operation\":1}]",
		FromUser:  "as23d",
		ChatID:    "asdf",
		Operation: network.SYNC_RESPONSE,
	}
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(mockSyncResponse)

}

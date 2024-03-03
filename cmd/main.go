package main

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
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
			From:      user1,
			To:        user2,
			Operation: c_model.SEND_MESSAGE,
		}
		message2 := c_model.Message{
			Id:        "2",
			Timestamp: 1633029443,
			Content:   "Hello, user1!",
			From:      user2,
			To:        user1,
			Operation: c_model.SEND_MESSAGE,
		}
		message3 := c_model.Message{
			Id:        "3",
			Timestamp: 1633029444,
			Content:   "Hello, user3!",
			From:      user1,
			To:        user3,
			Operation: c_model.SEND_MESSAGE,
		}
		message4 := c_model.Message{
			Id:        "4",
			Timestamp: 1633029445,
			Content:   "Hello, user4!",
			From:      user2,
			To:        user4,
			Operation: c_model.SEND_MESSAGE,
		}

		chats := p_model.NewNetworkChatMessages()
		chats.AddMessages([]c_model.Message{message1, message2, message3, message4})

		fmt.Printf("Chats: %v\n", chats.GetMessages())
	*/

	user2 := c_model.User{Username: "user2", UserId: "2"}
	user4 := c_model.User{Username: "user4", UserId: "4"}

	testMessage := c_model.Message{
		Id:        "8888",
		Timestamp: 1633029445,
		Content:   "Hello World!",
		From:      user2,
		To:        user4,
		Operation: c_model.TEST_MESSAGE,
	}

	p_model.GetPeerInstance()

	mockNetworkConnection := networkMockAdapter.GetMockConnection()
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(testMessage)
}

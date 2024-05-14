package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/stretchr/testify/assert"
	"testing"
)

// WIP

func setup() (*p_model.NetworkChatMessages, network.Message, network.Message, network.Message, network.Message) {
	ncm := p_model.NewNetworkChatMessages()

	message1 := network.Message{
		Id:              "1",
		Timestamp:       1633029442,
		Content:         "{\"message\": \"Hello, user2!\"}",
		SenderID:        "user1",
		ReceiverID:      "user2",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "user2.onion",
		ChatID:          "chat1",
		Operation:       network.SEND_MESSAGE,
	}
	message2 := network.Message{
		Id:              "2",
		Timestamp:       1633029443,
		Content:         "{\"message\": \"Hello, user1!\"}",
		SenderID:        "user2",
		ReceiverID:      "user1",
		SenderAddress:   "user2.onion",
		ReceiverAddress: "user1.onion",
		ChatID:          "chat1",
		Operation:       network.SEND_MESSAGE,
	}
	message3 := network.Message{
		Id:              "3",
		Timestamp:       1633029444,
		Content:         "{\"message\": \"Hello, user3!\"}",
		SenderID:        "user1",
		ReceiverID:      "user3",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "user3.onion",
		ChatID:          "chat2",
		Operation:       network.SEND_MESSAGE,
	}
	message4 := network.Message{
		Id:              "4",
		Timestamp:       1633029445,
		Content:         "{\"message\": \"Hello, user4!\"}",
		SenderID:        "user2",
		ReceiverID:      "user4",
		SenderAddress:   "user2.onion",
		ReceiverAddress: "user4.onion",
		ChatID:          "chat3",
		Operation:       network.SEND_MESSAGE,
	}

	return ncm, message1, message2, message3, message4
}

func TestAddMessage(t *testing.T) {
	ncm, message1, _, _, _ := setup()

	err := ncm.AddMessage(message1)
	assert.Nil(t, err)

	err = ncm.AddMessage(message1)
	assert.NotNil(t, err)
}

func TestGetMessages(t *testing.T) {
	ncm, message1, message2, _, _ := setup()

	_ = ncm.AddMessage(message1)
	_ = ncm.AddMessage(message2)

	messages := ncm.GetMessages()
	assert.Equal(t, 2, len(messages))
	assert.Contains(t, messages, message1)
	assert.Contains(t, messages, message2)
}

func TestGetMissingExternalMessages(t *testing.T) {
	ncm, message1, message2, _, _ := setup()

	_ = ncm.AddMessage(message1)

	missingExternal := ncm.GetMissingExternalMessages([]string{message2.Id})
	assert.Equal(t, 1, len(missingExternal))
	assert.Contains(t, missingExternal, message1)

	missingExternal = ncm.GetMissingExternalMessages([]string{message1.Id, message2.Id})
	assert.Equal(t, 0, len(missingExternal))
}

func TestGetMissingInternalMessages(t *testing.T) {
	ncm, message1, message2, message3, message4 := setup()

	_ = ncm.AddMessage(message1)
	_ = ncm.AddMessage(message3)

	missingInternal := ncm.GetMissingInternalMessages([]string{message4.Id})
	assert.Equal(t, 1, len(missingInternal))
	assert.Contains(t, missingInternal, message4)

	missingInternal = ncm.GetMissingInternalMessages([]string{message1.Id, message2.Id, message3.Id, message4.Id})
	assert.Equal(t, 2, len(missingInternal))
	assert.Contains(t, missingInternal, message2)
	assert.Contains(t, missingInternal, message4)
}

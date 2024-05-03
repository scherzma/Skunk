package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() (*p_model.NetworkChatMessages, network.Message, network.Message, network.Message, network.Message) {
	ncm := p_model.NewNetworkChatMessages()

	message1 := network.Message{
		Id:        "1",
		Timestamp: 1633029442,
		Content:   "Hello, user2!",
		FromUser:  "user1",
		ChatID:    "user2",
		Operation: network.SEND_MESSAGE,
	}
	message2 := network.Message{
		Id:        "2",
		Timestamp: 1633029443,
		Content:   "Hello, user1!",
		FromUser:  "user2",
		ChatID:    "user1",
		Operation: network.SEND_MESSAGE,
	}
	message3 := network.Message{
		Id:        "3",
		Timestamp: 1633029444,
		Content:   "Hello, user3!",
		FromUser:  "user1",
		ChatID:    "user3",
		Operation: network.SEND_MESSAGE,
	}
	message4 := network.Message{
		Id:        "4",
		Timestamp: 1633029445,
		Content:   "Hello, user4!",
		FromUser:  "user2",
		ChatID:    "user4",
		Operation: network.SEND_MESSAGE,
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

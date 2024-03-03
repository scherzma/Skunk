package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func setup() (*p_model.NetworkChatMessages, c_model.Message, c_model.Message, c_model.Message, c_model.Message) {
	ncm := p_model.NewNetworkChatMessages()

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

	missingExternal := ncm.GetMissingExternalMessages([]c_model.Message{message2})
	assert.Equal(t, 1, len(missingExternal))
	assert.Contains(t, missingExternal, message1)

	missingInternal := ncm.GetMissingInternalMessages([]c_model.Message{message1, message2})
	assert.Equal(t, 1, len(missingInternal))
}

func TestGetMissingInternalMessages(t *testing.T) {
	ncm, message1, message2, message3, message4 := setup()

	_ = ncm.AddMessage(message1)
	_ = ncm.AddMessage(message3)

	missingInternal := ncm.GetMissingInternalMessages([]c_model.Message{message4})

	assert.Equal(t, 1, len(missingInternal))
	assert.Contains(t, missingInternal, message4)

	missingInternal = ncm.GetMissingInternalMessages([]c_model.Message{message1, message2, message3, message4})
	assert.Equal(t, 2, len(missingInternal))
	assert.Contains(t, missingInternal, message4)
	assert.Contains(t, missingInternal, message2)
}

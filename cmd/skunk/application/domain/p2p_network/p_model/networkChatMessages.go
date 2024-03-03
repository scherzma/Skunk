package p_model

import (
	"errors"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
)

type NetworkChatMessages struct {
	chatMessagesMap map[string]c_model.Message
}

func NewNetworkChatMessages() *NetworkChatMessages {
	return &NetworkChatMessages{
		chatMessagesMap: make(map[string]c_model.Message),
	}
}

func (ncm *NetworkChatMessages) AddMessage(message c_model.Message) error {
	if _, exists := ncm.chatMessagesMap[message.Id]; exists {
		return errors.New("message already exists")
	}
	ncm.chatMessagesMap[message.Id] = message
	return nil
}

func (ncm *NetworkChatMessages) AddMessages(messages []c_model.Message) {
	for _, message := range messages {
		ncm.AddMessage(message)
	}
}

func (ncm *NetworkChatMessages) GetMessages() []c_model.Message {
	messages := make([]c_model.Message, 0, len(ncm.chatMessagesMap))
	for _, message := range ncm.chatMessagesMap {
		messages = append(messages, message)
	}
	return messages
}

func (ncm *NetworkChatMessages) GetMissingExternalMessages(inputMessages []c_model.Message) []c_model.Message {
	missingMessages := make([]c_model.Message, 0)
	inputMap := make(map[string]c_model.Message)

	// Convert input slice to map for efficient lookup
	for _, message := range inputMessages {
		inputMap[message.Id] = message
	}

	for _, message := range ncm.chatMessagesMap {
		if _, exists := inputMap[message.Id]; !exists {
			missingMessages = append(missingMessages, message)
		}
	}
	return missingMessages
}

func (ncm *NetworkChatMessages) GetMissingInternalMessages(inputMessages []c_model.Message) []c_model.Message {
	missingMessages := make([]c_model.Message, 0)
	ncmMap := ncm.chatMessagesMap

	for _, message := range inputMessages {
		if _, exists := ncmMap[message.Id]; !exists {
			missingMessages = append(missingMessages, message)
		}
	}
	return missingMessages
}

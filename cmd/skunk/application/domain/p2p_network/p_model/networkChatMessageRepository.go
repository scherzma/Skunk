package p_model

import (
	"errors"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type NetworkChatMessages struct {
	chatMessagesMap map[string]network.Message
}

func NewNetworkChatMessages() *NetworkChatMessages {
	return &NetworkChatMessages{
		chatMessagesMap: make(map[string]network.Message),
	}
}

func (ncm *NetworkChatMessages) AddMessage(message network.Message) error {
	if _, exists := ncm.chatMessagesMap[message.Id]; exists {
		return errors.New("message already exists")
	}
	ncm.chatMessagesMap[message.Id] = message
	return nil
}

func (ncm *NetworkChatMessages) GetMessages() []network.Message {
	messages := make([]network.Message, 0, len(ncm.chatMessagesMap))
	for _, message := range ncm.chatMessagesMap {
		messages = append(messages, message)
	}
	return messages
}

func (ncm *NetworkChatMessages) GetMissingExternalMessages(inputMessageIDs []string) []network.Message {
	missingMessages := make([]network.Message, 0)
	inputMap := make(map[string]bool)

	// Convert input slice to map for efficient lookup
	for _, id := range inputMessageIDs {
		inputMap[id] = true
	}

	for _, message := range ncm.chatMessagesMap {
		if !inputMap[message.Id] {
			missingMessages = append(missingMessages, message)
		}
	}
	return missingMessages
}

func (ncm *NetworkChatMessages) GetMissingInternalMessages(inputMessageIDs []string) []network.Message {
	missingMessages := make([]network.Message, 0)
	ncmMap := ncm.chatMessagesMap

	for _, id := range inputMessageIDs {
		if message, exists := ncmMap[id]; !exists {
			missingMessages = append(missingMessages, message)
		}
	}
	return missingMessages
}

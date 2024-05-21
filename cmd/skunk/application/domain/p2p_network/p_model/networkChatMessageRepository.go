// Package p_model provides data models for the p2p network.
package p_model

import (
	"errors"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// NetworkChatMessages represents a collection of chat messages in the network.
type NetworkChatMessages struct {
	chatName        string
	chatMessagesMap map[string]network.Message
}

// NewNetworkChatMessages creates a new instance of NetworkChatMessages.
func NewNetworkChatMessages() *NetworkChatMessages {
	return &NetworkChatMessages{
		chatMessagesMap: make(map[string]network.Message),
	}
}

// AddMessage adds a new message to the NetworkChatMessages.
// It returns an error if the message already exists.
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

// GetMissingExternalMessages returns the messages that are missing from the input message IDs.
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

// GetMissingInternalMessages returns the messages that are missing internally based on the input message IDs.
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

// GetMissingInternalMessageIDs returns the message IDs that are missing internally based on the input message IDs.
func (ncm *NetworkChatMessages) GetMissingInternalMessageIDs(inputMessageIDs []string) []string {
	missingMessageIDs := make([]string, 0)
	ncmMap := ncm.chatMessagesMap

	for _, id := range inputMessageIDs {
		if _, exists := ncmMap[id]; !exists {
			missingMessageIDs = append(missingMessageIDs, id)
		}
	}
	return missingMessageIDs
}

func (ncm *NetworkChatMessages) GetUsername() string {
	return "todo: implement me (Username)" // TODO implement me
}

func (ncm *NetworkChatMessages) GetChatName() string {
	return ncm.chatName
}

func (ncm *NetworkChatMessages) SetChatName(chatName string) {
	ncm.chatName = chatName
}

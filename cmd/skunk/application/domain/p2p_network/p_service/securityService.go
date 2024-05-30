package p_service

import (
	"encoding/json"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
)

type SecurityValidater interface {
	ValidateOutgoingMessage(message network.Message) bool
	ValidateIncomingMessage(message network.Message) bool
	ValidatePeer(peer string) bool
}

// SecurityContext is a service that provides security checks for the network
// It should be possible to implement all security checks in this service
// TODO: put a bit more thought into this
type SecurityContext struct {
	store           store.ChatInvitationStoragePort
	chatActionStore store.ChatActionStoragePort
	displayStorage  store.DisplayStoragePort
}

func NewSecurityContext(displayStorage store.DisplayStoragePort, store store.ChatInvitationStoragePort, chatActionStore store.ChatActionStoragePort) *SecurityContext {
	return &SecurityContext{
		store:           store,
		chatActionStore: chatActionStore,
		displayStorage:  displayStorage,
	}
}

func (s *SecurityContext) ValidateOutgoingMessage(message network.Message) bool {
	return true
}

func (s *SecurityContext) ValidateIncomingMessage(message network.Message) bool {
	switch message.Operation {
	case network.SEND_MESSAGE:
		return s.isMemberOfChat(message.SenderID, message.ChatID)
	case network.SYNC_REQUEST:
		return s.isMemberOfChat(message.SenderID, message.ChatID)
	case network.SYNC_RESPONSE:
		return s.validateSyncResponseMessages(message)
	case network.JOIN_CHAT:
		return s.hasValidInvitation(message.SenderID, message.ChatID)
	case network.LEAVE_CHAT:
		return s.isMemberOfChat(message.SenderID, message.ChatID)
	case network.INVITE_TO_CHAT:
		return true
	case network.SEND_FILE:
		return s.isMemberOfChat(message.SenderID, message.ChatID)
	case network.SET_USERNAME:
		return s.isMemberOfChat(message.SenderID, message.ChatID)
	case network.TEST_MESSAGE:
		return true
	case network.TEST_MESSAGE_2:
		return true
	default:
		return false
	}
}

func (s *SecurityContext) ValidatePeer(peer string) bool {
	return true
}

// Helper methods for security checks
func (s *SecurityContext) isMemberOfChat(peerID, chatID string) bool {
	members, err := s.displayStorage.GetUsersInChat(chatID)
	if err != nil {
		return false
	}

	for _, member := range members {
		if member.UserId == peerID {
			return true
		}
	}

	return false
}

func (s *SecurityContext) hasValidInvitation(peerID, chatID string) bool {
	invitations, err := s.store.GetInvitations(peerID)
	if err != nil {
		return false
	}

	for _, invitation := range invitations {
		if invitation == chatID {
			return true
		}
	}

	return false
}

func (s *SecurityContext) validateSyncResponseMessages(message network.Message) bool {
	var receivedMessages []network.Message
	err := json.Unmarshal([]byte(message.Content), &receivedMessages)
	if err != nil {
		return false
	}

	for _, msg := range receivedMessages {
		if msg.SenderID != message.SenderID {
			return false
		}
	}

	return true
}

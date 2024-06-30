package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
)

// InviteToChatHandler handles invitations for a peer to join a chat
type InviteToChatHandler struct {
	userChatLogic         chat.ChatLogic
	chatInvitationStorage store.ChatInvitationStoragePort
}

// NewInviteToChatHandler creates a new InviteToChatHandler
func NewInviteToChatHandler(userChatLogic chat.ChatLogic, chatInvitationStorage store.ChatInvitationStoragePort) *InviteToChatHandler {
	return &InviteToChatHandler{
		userChatLogic:         userChatLogic,
		chatInvitationStorage: chatInvitationStorage,
	}
}

// HandleMessage processes the chat invitation message
func (i *InviteToChatHandler) HandleMessage(message network.Message) error {

	// Message structure:
	/*
		{
			"chatId": "chat_id",
			"chatName": "chat_name",
			"peers": [
				{
					"address": "peer_address",
					"publicKey": "peer_public_key"
				},
				...
			]
		}
	*/

	var content struct {
		ChatID   string                   `json:"chatId"`
		ChatName string                   `json:"chatName"`
		Peers    []store.PublicKeyAddress `json:"peers"`
	}

	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	// Store the chat invitation details
	err = i.chatInvitationStorage.InvitedToChat(message.Id, content.Peers)
	if err != nil {
		return err
	}

	// Extract addresses from peers
	peerAddresses := make([]string, len(content.Peers))
	for idx, peer := range content.Peers {
		peerAddresses[idx] = peer.Address
	}

	// Notify the chat logic of the received invitation
	i.userChatLogic.ReceiveChatInvitation(message.SenderID, content.ChatID, content.ChatName, peerAddresses)

	return nil
}

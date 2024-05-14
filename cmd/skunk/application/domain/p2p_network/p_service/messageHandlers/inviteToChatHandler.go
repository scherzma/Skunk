package messageHandlers

import (
	"encoding/json"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
)

// This Peer gets invited to a chat
type inviteToChatHandler struct {
	userChatLogic         chat.ChatLogic
	chatInvitationStorage store.ChatInvitationStoragePort
}

func NewInviteToChatHandler(userChatLogic chat.ChatLogic, chatInvitationStorage store.ChatInvitationStoragePort) *inviteToChatHandler {
	return &inviteToChatHandler{
		userChatLogic:         userChatLogic,
		chatInvitationStorage: chatInvitationStorage,
	}
}

func (i *inviteToChatHandler) HandleMessage(message network.Message) error {

	// Structure of the message:
	/*
		{
			"chatId": "asdf",
			"chatName": "asdf",
			"peers": [
				"asdf",
				"asdf"
			]
		}
	*/

	var content struct {
		ChatID   string `json:"chatId"`
		ChatName string `json:"chatName"`
		Peers    []string
	}

	err := json.Unmarshal([]byte(message.Content), &content)
	if err != nil {
		fmt.Println("Error unmarshalling message content")
		return err
	}

	err = i.chatInvitationStorage.InvitedToChat(message.Id, []store.PublicKeyAddress{{Address: content.ChatName}})
	if err != nil {
		return err
	}

	i.userChatLogic.ReceiveChatInvitation(message.SenderID, content.ChatID, content.ChatName, content.Peers)

	return nil
}

package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
)

// A Peer joins a chat
type joinChatHandler struct {
	userChatLogic     chat.ChatLogic
	chatActionStorage store.ChatActionStoragePort
}

func NewJoinChatHandler(userChatLogic chat.ChatLogic, chatActionStorage store.ChatActionStoragePort) *joinChatHandler {
	return &joinChatHandler{
		userChatLogic:     userChatLogic,
		chatActionStorage: chatActionStorage,
	}
}

func (j *joinChatHandler) HandleMessage(message network.Message) error {

	// Update chat invitation storage
	err := j.chatActionStorage.PeerJoinedChat(message.SenderID, message.ChatID)
	if err != nil {
		fmt.Println("Error updating chat invitation storage")
		return err
	}

	// Handle peer joining the chat
	j.userChatLogic.PeerJoinsChat(message.SenderID, message.ChatID)

	return nil
}

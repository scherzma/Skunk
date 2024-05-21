package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
)

// A Peer leaves a chat
type leaveChatHandler struct {
	userChatLogic     chat.ChatLogic
	chatActionStorage store.ChatActionStoragePort
}

func NewLeaveChatHandler(userChatLogic chat.ChatLogic, chatActionStorage store.ChatActionStoragePort) *leaveChatHandler {
	return &leaveChatHandler{
		userChatLogic:     userChatLogic,
		chatActionStorage: chatActionStorage,
	}
}

func (l *leaveChatHandler) HandleMessage(message network.Message) error {

	// Update chat invitation storage and remove the peer from the chat
	err := l.chatActionStorage.PeerLeftChat(message.SenderID, message.ChatID)
	if err != nil {
		fmt.Println("Error updating chat invitation storage")
		return err
	}

	// Handle peer leaving the chat
	l.userChatLogic.PeerLeavesChat(message.SenderID, message.ChatID)

	return nil
}

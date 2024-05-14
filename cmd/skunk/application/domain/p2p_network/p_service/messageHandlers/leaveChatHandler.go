package messageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// A Peer leaves a chat
type LeaveChatHandler struct {
	userChatLogic chat.ChatLogic
}

func (l *LeaveChatHandler) HandleMessage(message network.Message) error {

	l.userChatLogic.PeerLeavesChat(message.SenderID, message.ChatID)
	//TODO make the necessary changes to the chat (SQLite (with interface of course))

	return nil
}

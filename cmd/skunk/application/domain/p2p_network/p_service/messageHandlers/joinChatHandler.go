package messageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// A Peer joins a chat
type JoinChatHandler struct {
	userChatLogic chat.ChatLogic
}

func (j *JoinChatHandler) HandleMessage(message network.Message) error {

	j.userChatLogic.PeerJoinsChat(message.SenderID, message.ChatID)
	//TODO make the necessary changes to the chat (SQLite (with interface of course))

	return nil
}

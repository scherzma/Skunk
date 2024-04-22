package messageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type CreateChatHandler struct {
}

func (handler *CreateChatHandler) HandleMessage(message network.Message) error {
	//TODO implement

	// Is this needed? In which case would i send a message to someone else about a created chat?
	return nil
}

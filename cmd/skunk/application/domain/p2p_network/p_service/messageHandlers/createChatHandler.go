package messageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type CreateChatHandler struct {
}

func (handler *CreateChatHandler) HandleMessage(message network.Message) error {
	//TODO implement
	return nil
}

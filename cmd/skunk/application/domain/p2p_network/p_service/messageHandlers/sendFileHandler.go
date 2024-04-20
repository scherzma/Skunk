package messageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type SendFileHandler struct{}

func (s *SendFileHandler) HandleMessage(message network.Message) error {
	//TODO implement
	return nil
}

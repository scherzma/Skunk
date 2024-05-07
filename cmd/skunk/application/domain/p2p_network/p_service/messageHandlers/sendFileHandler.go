package messageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// SendFileHandler handles the "SendFile" message operation.
type SendFileHandler struct{}

func (s *SendFileHandler) HandleMessage(message network.Message) error {
	//TODO implement
	return nil
}

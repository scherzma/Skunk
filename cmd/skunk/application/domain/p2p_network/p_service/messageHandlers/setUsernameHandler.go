package messageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// SetUsernameHandler handles the "SetUsername" message operation.
type SetUsernameHandler struct{}

func (s *SetUsernameHandler) HandleMessage(message network.Message) error {
	//TODO implement
	return nil
}

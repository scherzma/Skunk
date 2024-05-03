package frontendMessageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"
)

type SetUsernameHandler struct{}

func (s *SetUsernameHandler) HandleMessage(message frontend.FrontendMessage) error {
	// TODO: Implement
	return nil
}

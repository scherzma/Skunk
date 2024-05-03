package frontendMessageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"
)

type SendFileHandler struct{}

func (s *SendFileHandler) HandleMessage(message frontend.FrontendMessage) error {
	// TODO: Implement
	return nil
}

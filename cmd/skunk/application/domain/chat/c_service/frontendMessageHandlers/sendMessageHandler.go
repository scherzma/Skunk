package frontendMessageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"
)

type SendMessageHandler struct{}

func (s *SendMessageHandler) HandleMessage(message frontend.FrontendMessage) error {
	// TODO: Implement
	return nil
}

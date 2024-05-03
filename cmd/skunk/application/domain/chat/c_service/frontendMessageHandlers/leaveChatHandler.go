package frontendMessageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"
)

type LeaveChatHandler struct{}

func (l *LeaveChatHandler) HandleMessage(message frontend.FrontendMessage) error {
	// TODO: Implement
	return nil
}

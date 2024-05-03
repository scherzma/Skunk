package frontendMessageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"

type FrontendMessageHandler interface {
	HandleMessage(message frontend.FrontendMessage) error
}

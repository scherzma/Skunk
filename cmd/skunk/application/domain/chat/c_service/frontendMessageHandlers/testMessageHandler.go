package frontendMessageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/frontend"
)

type TestMessageHandler struct{}

func (t *TestMessageHandler) HandleMessage(message frontend.FrontendMessage) error {
	fmt.Println("TestMessageHandler: ", message)
	return nil
}

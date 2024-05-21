package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// TestMessageHandler handles the "TestMessage" message operation.
type TestMessageHandler struct {
}

func (t *TestMessageHandler) HandleMessage(message network.Message) error {

	fmt.Println("TestMessageHandler: ", message.Content)
	return nil
}

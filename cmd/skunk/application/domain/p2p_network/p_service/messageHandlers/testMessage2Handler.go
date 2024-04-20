package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type TestMessageHandler2 struct {
}

func (t *TestMessageHandler2) HandleMessage(message network.Message) error {
	fmt.Println("TestMessageHandler_2: ", message.Content)
	return nil
}

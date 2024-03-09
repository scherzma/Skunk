package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
)

type TestMessageHandler2 struct {
}

func (t *TestMessageHandler2) HandleMessage(message c_model.Message) error {
	fmt.Println("TestMessageHandler_2: ", message.Content)
	return nil
}

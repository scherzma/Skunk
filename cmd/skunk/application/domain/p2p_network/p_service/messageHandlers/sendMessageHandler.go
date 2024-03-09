package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
)

type SendMessageHandler struct{}

func (s *SendMessageHandler) HandleMessage(message c_model.Message) error {
	fmt.Println("HandleMessage: ", message.Content)
	return nil
}

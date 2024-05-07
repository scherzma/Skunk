package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// SendMessageHandler handles the "SendMessage" message operation.
type SendMessageHandler struct{}

func (s *SendMessageHandler) HandleMessage(message network.Message) error {
	fmt.Println("HandleMessage: ", message.Content)
	// TODO define interface for ChatLogic with Pub/Sub

	return nil
}

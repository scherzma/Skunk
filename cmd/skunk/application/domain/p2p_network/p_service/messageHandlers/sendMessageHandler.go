package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type SendMessageHandler struct{}

func (s *SendMessageHandler) HandleMessage(message network.Message) error {
	fmt.Println("HandleMessage: ", message.Content)

	// peer := network.GetPeer(message.Sender)

	return nil
}

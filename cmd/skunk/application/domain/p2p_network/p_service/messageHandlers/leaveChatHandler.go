package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type LeaveChatHandler struct{}

func (l *LeaveChatHandler) HandleMessage(message network.Message) error {
	//TODO implement
	fmt.Println("LeaveChatHandler")
	return nil
}

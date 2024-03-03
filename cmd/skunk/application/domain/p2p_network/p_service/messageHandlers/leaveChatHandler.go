package messageHandlers

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
)

type LeaveChatHandler struct{}

func (l *LeaveChatHandler) HandleMessage(message c_model.Message) error {
	//TODO implement
	fmt.Println("LeaveChatHandler")
	return nil
}

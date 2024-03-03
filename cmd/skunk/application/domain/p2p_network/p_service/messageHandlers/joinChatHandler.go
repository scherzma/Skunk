package messageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"

type JoinChatHandler struct{}

func (j *JoinChatHandler) HandleMessage(message c_model.Message) error {
	//TODO implement
	return nil
}

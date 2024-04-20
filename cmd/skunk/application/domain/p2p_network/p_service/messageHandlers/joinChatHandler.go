package messageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/port/network"

type JoinChatHandler struct{}

func (j *JoinChatHandler) HandleMessage(message network.Message) error {
	//TODO implement
	return nil
}

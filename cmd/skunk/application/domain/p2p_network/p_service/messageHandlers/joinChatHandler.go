package messageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/port/network"

// JoinChatHandler handles the "JoinChat" message operation.
type JoinChatHandler struct{}

func (j *JoinChatHandler) HandleMessage(message network.Message) error {
	//TODO implement
	return nil
}

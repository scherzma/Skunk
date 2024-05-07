package messageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/port/network"

// InviteToChatHandler handles the "InviteToChat" message operation.
type InviteToChatHandler struct{}

func (i *InviteToChatHandler) HandleMessage(message network.Message) error {
	//TODO implement
	return nil
}

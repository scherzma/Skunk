package messageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/port/network"

// MessageHandler is an interface that defines the contract for handling network messages.
// Types that implement this interface can be used to process specific message operations.
type MessageHandler interface {
	// HandleMessage processes the received network message.
	// It takes a network.Message as input and returns an error if the handling fails.
	HandleMessage(message network.Message) error
}

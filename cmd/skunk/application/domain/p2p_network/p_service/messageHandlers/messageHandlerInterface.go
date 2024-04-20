package messageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/port/network"

type MessageHandler interface {
	HandleMessage(message network.Message) error
}

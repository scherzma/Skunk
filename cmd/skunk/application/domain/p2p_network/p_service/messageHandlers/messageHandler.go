package messageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"

type MessageHandler interface {
	HandleMessage(message c_model.Message) error
}

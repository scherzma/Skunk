package port

import "github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"

type StoreMessageQueue interface {
	StoreMessageQueue(messageQueue p_model.MessageQueue) error
	UpdateMessageQueue(messageQueue p_model.MessageQueue) error
	RetriveMessageQueue() (string, error)
}

package network

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
)

type NetworkObserver interface {
	Notify(message c_model.Message) error
}

type NetworkConnection interface {
	SubscribeToNetwork(observer NetworkObserver) error
	UnsubscribeFromNetwork(observer NetworkObserver) error
	SendMessageToNetworkPeer(address string, message c_model.Message) error
}

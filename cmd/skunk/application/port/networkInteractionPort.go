package port

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
)

type NetworkConnection interface {
	SubscribeToNetwork(observer p_model.NetworkObserver) error
	UnsubscribeFromNetwork(observer p_model.NetworkObserver) error
	SendMessageToNetworkPeer(address string, message c_model.Message) error
}

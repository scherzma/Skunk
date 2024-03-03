package p_model

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port"
)

type NetworkPeer struct {
	Id                   string
	NetworkPeerPublicKey string
	NetworkPeerAddress   string
}

func (n *NetworkPeer) SendMessage(message c_model.Message) error {
	var con port.NetworkConnection = &networkMockAdapter.MockConnection{}
	con.SendMessageToNetworkPeer(n.NetworkPeerAddress, message)
	return nil
}

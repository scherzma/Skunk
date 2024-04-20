package p_model

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type NetworkPeer struct {
	Id                   string
	NetworkPeerPublicKey string
	NetworkPeerAddress   string
	Connection           network.NetworkConnection
}

func (n *NetworkPeer) SendMessage(message network.Message) error {
	n.Connection.SendMessageToNetworkPeer(n.NetworkPeerAddress, message)
	return nil
}

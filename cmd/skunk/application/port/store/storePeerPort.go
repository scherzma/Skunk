package store

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
)

type StorePeer interface {
	StorePeer(peer messageHandlers.Peer)
	RetrivePeers() ([]messageHandlers.Peer, error)
	RetrivePeer(peerId string) (messageHandlers.Peer, error)
}

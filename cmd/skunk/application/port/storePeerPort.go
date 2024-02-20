package port

import "github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"

type StorePeer interface {
	StorePeer(peer p_model.Peer)
	RetrivePeers() ([]p_model.Peer, error)
	RetrivePeer(peerId string) (p_model.Peer, error)
}

package test

import (
	"testing"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/util"
	"github.com/stretchr/testify/assert"
)

func TestNetworkAdapter(t *testing.T) {
	testMessage := network.Message{
		Id:        util.UUID(),
		Timestamp: util.CurrentTimeMillis(),
		Content:   "Hello World!",
		FromUser:  "Alice",
		ChatID:    "1",
		Operation: network.TEST_MESSAGE,
	}

	peerInstance := messageHandlers.GetPeerInstance()

	networkConnection := networkAdapter.NewAdapter()
	peerInstance.AddNetworkConnection(networkConnection)

	err := peerInstance.SendMessageToNetworkPeer("", testMessage)
	assert.NoError(t, err)

	err = peerInstance.SendMessageToNetworkPeer("", testMessage)
	assert.NoError(t, err)

	peerInstance.RemoveNetworkConnection(networkConnection)
}

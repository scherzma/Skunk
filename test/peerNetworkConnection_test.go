package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGetPeerInstance tests the GetPeerInstance function of the messageHandlers package.
// It verifies that the function returns the same instance when called multiple times.
func TestGetPeerInstance(t *testing.T) {
	peer1 := p_service.GetPeerInstance()
	peer2 := p_service.GetPeerInstance()

	if peer1 != peer2 {
		t.Errorf("GetPeerInstance() failed, expected same instance, got different instances")
	}
}

// TestNotify tests the Notify method of the Peer struct.
// It creates a test message and sends it to the peer using the Notify method.
// It asserts that no error is returned.
func TestNotify(t *testing.T) {
	peer := p_service.GetPeerInstance()

	testMessage := network.Message{
		Id:              "8888",
		Timestamp:       1633029445,
		Content:         "Hello World!",
		SenderID:        "user1",
		ReceiverID:      "user2",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "user2.onion",
		ChatID:          "chat1",
		Operation:       network.TEST_MESSAGE,
	}

	err := peer.Notify(testMessage)

	assert.NoError(t, err, "Notify() failed, expected nil, got error")

}

// TestSubscribeAndUnsubscribeToNetwork tests the SubscribeToNetwork and UnsubscribeFromNetwork
// methods of the NetworkConnection interface.
// It uses a mock connection and a peer instance to subscribe and unsubscribe the peer
// to/from the network.
// It asserts that no error is returned in both cases.
func TestSubscribeAndUnsubscribeToNetwork(t *testing.T) {
	mockConnection := networkMockAdapter.GetMockConnection()
	peer := p_service.GetPeerInstance()

	err := mockConnection.SubscribeToNetwork(peer)
	assert.NoError(t, err, "SubscribeToNetwork() failed, expected nil, got error")

	err = mockConnection.UnsubscribeFromNetwork()
	assert.NoError(t, err, "UnsubscribeFromNetwork() failed, expected nil, got error")
}

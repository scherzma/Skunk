package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"testing"
)

func TestGetPeerInstance(t *testing.T) {
	peer1 := p_model.GetPeerInstance()
	peer2 := p_model.GetPeerInstance()

	if peer1 != peer2 {
		t.Errorf("GetPeerInstance() failed, expected same instance, got different instances")
	}
}

func TestNotify(t *testing.T) {
	peer := p_model.GetPeerInstance()

	testMessage := network.Message{
		Id:        "8888",
		Timestamp: 1633029445,
		Content:   "Hello World!",
		Operation: network.TEST_MESSAGE,
	}

	err := peer.Notify(testMessage)

	if err != nil {
		t.Errorf("Notify() failed, expected nil, got error: %v", err)
	}
}

func TestSubscribeAndUnsubscribeToNetwork(t *testing.T) {
	mockConnection := networkMockAdapter.GetMockConnection()
	peer := p_model.GetPeerInstance()

	err := mockConnection.SubscribeToNetwork(peer)
	if err != nil {
		t.Errorf("SubscribeToNetwork() failed, expected nil, got error: %v", err)
	}

	err = mockConnection.UnsubscribeFromNetwork(peer)
	if err != nil {
		t.Errorf("UnsubscribeFromNetwork() failed, expected nil, got error: %v", err)
	}
}

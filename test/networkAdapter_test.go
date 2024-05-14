package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"testing"
)

type MockPeer struct {
	network network.NetworkConnection
}

func (p *MockPeer) SubscribeToNetwork(network network.NetworkConnection) error {
	p.network = network
	p.network.SubscribeToNetwork(p)
	return nil
}

func (p *MockPeer) RemoveNetworkConnection(network network.NetworkConnection) {
	p.network.UnsubscribeFromNetwork(p)
	p.network = nil
}

func (p *MockPeer) Notify(message network.Message) error {
	return nil
}

func (p *MockPeer) SendMessageToNetworkPeer(address string, message network.Message) error {
	p.network.SendMessageToNetworkPeer(address, message)
	return nil
}

func TestNetworkAdapter(t *testing.T) {

	// Create a mock network connection
	testMessage := network.Message{
		Id:              "msg1",
		Timestamp:       1620000000,
		Content:         "{\"message\": \"Hey everyone!\"}",
		SenderID:        "user1",
		ReceiverID:      "user2",
		SenderAddress:   "user1.onion",
		ReceiverAddress: "user2.onion",
		ChatID:          "chat1",
		Operation:       network.SEND_MESSAGE,
	}

	peer := MockPeer{
		network: &networkMockAdapter.MockConnection{},
	}

	mockNetworkConnection := networkMockAdapter.GetMockConnection()
	peer.SubscribeToNetwork(mockNetworkConnection)

	peer.SendMessageToNetworkPeer("addressResponse", testMessage)
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(testMessage)

}

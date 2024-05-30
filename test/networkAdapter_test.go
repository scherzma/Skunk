package test

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"testing"
)

// Implement a mockPeer to test the networkAdapter, not the peer.
type MockPeer struct {
	network          network.NetworkConnection
	receivedMessages []network.Message
}

func (p *MockPeer) SubscribeToNetwork(network network.NetworkConnection) error {
	p.network = network
	p.network.SubscribeToNetwork(p)
	return nil
}

func (p *MockPeer) RemoveNetworkConnection(network network.NetworkConnection) {
	p.network.UnsubscribeFromNetwork()
	p.network = nil
}

func (p *MockPeer) Notify(message network.Message) error {
	p.receivedMessages = append(p.receivedMessages, message)
	return nil
}

func (p *MockPeer) SendMessageToNetworkPeer(address string, message network.Message) error {
	p.network.SendMessageToNetworkPeer(message)
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

	peer := &MockPeer{
		network: &networkMockAdapter.MockConnection{},
	}

	mockNetworkConnection := networkMockAdapter.GetMockConnection()
	err := peer.SubscribeToNetwork(mockNetworkConnection)
	if err != nil {
		t.Fatalf("Failed to subscribe to network: %v", err)
	}

	peer.SendMessageToNetworkPeer("addressResponse", testMessage)
	mockNetworkConnection.SendMockNetworkMessageToSubscribers(testMessage)

	if len(peer.receivedMessages) != 1 {
		t.Fatalf("Expected 1 received message, got %d", len(peer.receivedMessages))
	}

	receivedMessage := peer.receivedMessages[0]
	if receivedMessage.Id != testMessage.Id {
		t.Errorf("Expected message ID %s, got %s", testMessage.Id, receivedMessage.Id)
	}
	if receivedMessage.Content != testMessage.Content {
		t.Errorf("Expected message content %s, got %s", testMessage.Content, receivedMessage.Content)
	}
	if receivedMessage.SenderID != testMessage.SenderID {
		t.Errorf("Expected sender ID %s, got %s", testMessage.SenderID, receivedMessage.SenderID)
	}
	if receivedMessage.ReceiverID != testMessage.ReceiverID {
		t.Errorf("Expected receiver ID %s, got %s", testMessage.ReceiverID, receivedMessage.ReceiverID)
	}
	if receivedMessage.ChatID != testMessage.ChatID {
		t.Errorf("Expected chat ID %s, got %s", testMessage.ChatID, receivedMessage.ChatID)
	}
}

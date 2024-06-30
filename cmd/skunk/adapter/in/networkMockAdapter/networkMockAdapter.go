package networkMockAdapter

import (
	"fmt"
	"sync"

	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// NetworkMockAdapter is a mock adapter for the network
// It implements the NetworkConnection interface
var (
	mockConnection *MockConnection // singleton for testing purposes
	once           sync.Once
)

type MockConnection struct {
	subscriber network.NetworkObserver
	LastSent   network.Message
}

func GetMockConnection() *MockConnection {
	once.Do(func() {
		mockConnection = &MockConnection{}
	})
	return mockConnection
}

// SubscribeToNetwork is a mock function for the network
func (m *MockConnection) SubscribeToNetwork(observer network.NetworkObserver) error {
	if m.subscriber != nil {
		return fmt.Errorf("network adapter is already connected to observer: %v", observer)
	}

	m.subscriber = observer
	return nil
}

// UnsubscribeFromNetwork is a mock function for the network
func (m *MockConnection) UnsubscribeFromNetwork() error {
	if m.subscriber == nil {
		return fmt.Errorf("can't unsubscribe from nil")
	}

	m.subscriber = nil
	return nil
}

// SendMessageToNetworkPeer is a mock function for the network
func (m *MockConnection) SendMessageToNetworkPeer(message network.Message) error {
	fmt.Println("Sending message to: " + message.ReceiverAddress)
	fmt.Println("Message: ", message)
	m.LastSent = message
	return nil
}

// SendMockNetworkMessageToSubscribers is a mock function for the network
func (m *MockConnection) SendMockNetworkMessageToSubscribers(message network.Message) error {
	m.subscriber.Notify(message)
	return nil
}

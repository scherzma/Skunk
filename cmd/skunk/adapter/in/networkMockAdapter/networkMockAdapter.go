package networkMockAdapter

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"sync"
)

// NetworkMockAdapter is a mock adapter for the network
// It implements the NetworkConnection interface

var (
	mockConnection *MockConnection // singleton for testing purposes
	once           sync.Once
)

type MockConnection struct {
	subscribers []network.NetworkObserver
}

func GetMockConnection() *MockConnection {
	once.Do(func() {
		mockConnection = &MockConnection{}
	})
	return mockConnection
}

// SubscribeToNetwork is a mock function for the network
func (m *MockConnection) SubscribeToNetwork(observer network.NetworkObserver) error {
	m.subscribers = append(m.subscribers, observer)
	return nil
}

// UnsubscribeFromNetwork is a mock function for the network
func (m *MockConnection) UnsubscribeFromNetwork(observer network.NetworkObserver) error {
	for i, sub := range m.subscribers {
		if sub == observer {
			m.subscribers = append(m.subscribers[:i], m.subscribers[i+1:]...)
			break
		}
	}
	return nil
}

// SendMessageToNetworkPeer is a mock function for the network
func (m *MockConnection) SendMessageToNetworkPeer(address string, message network.Message) error {
	fmt.Println("Sending message to: " + address)
	fmt.Println(message.Content)
	return nil
}

// SendMockNetworkMessageToSubscribers is a mock function for the network
func (m *MockConnection) SendMockNetworkMessageToSubscribers(message network.Message) error {
	for _, sub := range m.subscribers {
		sub.Notify(message)
	}
	return nil
}

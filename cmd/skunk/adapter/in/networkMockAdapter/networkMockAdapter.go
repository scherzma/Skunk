package networkMockAdapter

import (
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
)

// NetworkMockAdapter is a mock adapter for the network
// It implements the NetworkConnection interface

type MockConnection struct {
	subscribers []p_model.NetworkObserver
}

// SubscribeToNetwork is a mock function for the network
func (m *MockConnection) SubscribeToNetwork(observer p_model.NetworkObserver) error {
	m.subscribers = append(m.subscribers, observer)
	return nil
}

// UnsubscribeFromNetwork is a mock function for the network
func (m *MockConnection) UnsubscribeFromNetwork(observer p_model.NetworkObserver) error {
	for i, sub := range m.subscribers {
		if sub == observer {
			m.subscribers = append(m.subscribers[:i], m.subscribers[i+1:]...)
			break
		}
	}
	return nil
}

// SendMessageToNetworkPeer is a mock function for the network
func (m *MockConnection) SendMessageToNetworkPeer(address string, message c_model.Message) error {
	fmt.Println("Sending message to: " + address)
	fmt.Println(message.Content)
	return nil
}

// SendMockNetworkMessageToSubscribers is a mock function for the network
func (m *MockConnection) SendMockNetworkMessageToSubscribers(message c_model.Message) error {
	for _, sub := range m.subscribers {
		sub.Notify(message)
	}
	return nil
}

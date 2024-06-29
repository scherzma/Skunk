package messageHandlers

import (
	"errors"
	"fmt"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

// This has many problems right now.
// For example is it only able to handle one connection. That means, while the peer is able to handle more, this is a bottleneck.
// One could create multiple MessageSenders, but this is probably a bad idea.
// Because, if done so, everything would need multiple messageSenders, or swap them for possibly ever message.
// Still, being able to handle multiple networkConnections is not a requirement for the first release, and while many parts of the architecture support it, this will suffice.
type MessageSender struct {
	networkConnection network.NetworkConnection
	securityContext   p_service.SecurityValidater
}

func NewMessageSender(securityContext p_service.SecurityValidater) *MessageSender {
	return &MessageSender{
		securityContext:   securityContext,
		networkConnection: nil,
	}
}

func (m *MessageSender) SendMessage(message network.Message) error {
	if m.networkConnection == nil {
		return fmt.Errorf("no network connection is set")
	}

	if !m.securityContext.ValidateOutgoingMessage(message) {
		return errors.New("invalid message")
	}

	err := m.networkConnection.SendMessageToNetworkPeer(message)
	if err != nil {
		return err
	}

	return nil
}

func (m *MessageSender) SetNetworkConnection(networkConnection network.NetworkConnection) {
	m.networkConnection = networkConnection
}

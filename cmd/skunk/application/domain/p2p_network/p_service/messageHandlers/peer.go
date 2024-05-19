// Package messageHandlers provides message handling for the p2p network.
package messageHandlers

import (
	"errors"
    "fmt"
	"sync"

	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

var (
	peerInstance *Peer
	once         sync.Once
)

// Peer represents a node in the p2p network. It manages network connections,
// message handlers, and security context. The Peer struct is a singleton instance
// that is lazily initialized using the sync.Once mechanism.
type Peer struct {
	Chats           p_model.NetworkChats
	Address         string
	handlers        map[network.OperationType]MessageHandler
	connections     []network.NetworkConnection
	securityContext p_service.SecurityContext
}

func GetPeerInstance() *Peer {
	once.Do(func() {
		handlers := map[network.OperationType]MessageHandler{
			network.SEND_MESSAGE:  &SendMessageHandler{},
			network.SYNC_REQUEST:  &SyncRequestHandler{},
			network.SYNC_RESPONSE: &SyncResponseHandler{},
			// CREATE_CHAT missing
			network.JOIN_CHAT:      &JoinChatHandler{},
			network.LEAVE_CHAT:     &LeaveChatHandler{},
			network.INVITE_TO_CHAT: &InviteToChatHandler{},
			network.SEND_FILE:      &SendFileHandler{},
			network.SET_USERNAME:   &SetUsernameHandler{},
			// USER_OFFLINE missing
			network.NETWORK_ONLINE: &NetworkOnlineHandler{},
			network.TEST_MESSAGE:   &TestMessageHandler{},
			network.TEST_MESSAGE_2: &TestMessageHandler2{},
		}

		peerInstance = &Peer{
			Chats:           p_model.NetworkChats{},
			handlers:        handlers,
			connections:     []network.NetworkConnection{},
			securityContext: p_service.SecurityContext{},
		}
	})

	return peerInstance
}

// AddNetworkConnection adds a new network connection to the Peer instance
// and subscribes the Peer to the network events.
func (p *Peer) AddNetworkConnection(connection network.NetworkConnection) {
	p.connections = append(p.connections, connection)
	connection.SubscribeToNetwork(p)
}

// RemoveNetworkConnection removes a network connection from the Peer instance.
func (p *Peer) RemoveNetworkConnection(connection network.NetworkConnection) {
	for i, c := range p.connections {
		if c == connection {
            err := c.UnsubscribeFromNetwork()
            if err != nil {
                fmt.Println(err)
            }
			p.connections = append(p.connections[:i], p.connections[i+1:]...)
			break
		}
	}
}

// Notify handles incoming network messages. It validates the message using the
// security context and routes it to the appropriate message handler based on
// the message operation type. If the message is invalid or the operation type
// is not supported, an error is returned.
func (p *Peer) Notify(message network.Message) error {
	if handler, exists := p.handlers[message.Operation]; exists {
		if !p.securityContext.ValidateIncomingMessage(message) {
			return errors.New("invalid message")
		}
		return handler.HandleMessage(message)
	}
	return errors.New("invalid message operation")
}

// SendMessageToNetworkPeer sends a message to a network peer.
func (p *Peer) SendMessageToNetworkPeer(address string, message network.Message) error {
	if !p.securityContext.ValidateOutgoingMessage(message) {
		return errors.New("invalid message")
	}

	for _, connection := range p.connections {
		if err := connection.SendMessageToNetworkPeer(address, message); err != nil {
			return err
		}
	}
	return nil
}

package messageHandlers

import (
	"errors"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"sync"
)

var (
	peerInstance *Peer
	once         sync.Once
)

type Peer struct {
	Chats       p_model.NetworkChats
	handlers    map[network.OperationType]MessageHandler
	connections []network.NetworkConnection
}

func GetPeerInstance() *Peer {
	once.Do(func() {
		handlers := map[network.OperationType]MessageHandler{
			network.JOIN_CHAT:      &JoinChatHandler{},
			network.SEND_FILE:      &SendFileHandler{},
			network.SYNC_REQUEST:   &SyncRequestHandler{},
			network.SYNC_RESPONSE:  &SyncResponseHandler{},
			network.SET_USERNAME:   &SetUsernameHandler{},
			network.SEND_MESSAGE:   &SendMessageHandler{},
			network.CREATE_CHAT:    &CreateChatHandler{},
			network.INVITE_TO_CHAT: &InviteToChatHandler{},
			network.LEAVE_CHAT:     &LeaveChatHandler{},
			network.TEST_MESSAGE:   &TestMessageHandler{},
			network.TEST_MESSAGE_2: &TestMessageHandler2{},
		}

		peerInstance = &Peer{
			Chats:       p_model.NetworkChats{},
			handlers:    handlers,
			connections: []network.NetworkConnection{},
		}
	})

	return peerInstance
}

func (p *Peer) AddNetworkConnection(connection network.NetworkConnection) {
	p.connections = append(p.connections, connection)
	connection.SubscribeToNetwork(p)
}

func (p *Peer) RemoveNetworkConnection(connection network.NetworkConnection) {
	for i, c := range p.connections {
		if c == connection {
			p.connections = append(p.connections[:i], p.connections[i+1:]...)
			break
		}
	}
}

func (p *Peer) Notify(message network.Message) error {
	if handler, exists := p.handlers[message.Operation]; exists {
		return handler.HandleMessage(message) // some form of authentication should be done here
	}
	return errors.New("invalid message operation")
}

func (p *Peer) SendMessageToNetworkPeer(address string, message network.Message) error {

	if !p_service.ValidateMessage(message) {
		return errors.New("invalid message")
	}

	for _, connection := range p.connections {
		if err := connection.SendMessageToNetworkPeer(address, message); err != nil {
			return err
		}
	}
	return nil
}
package p_model

import (
	"errors"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"sync"
)

var (
	peerInstance *Peer
	once         sync.Once
)

type Peer struct {
	Chats       networkChats
	handlers    map[network.OperationType]messageHandlers.MessageHandler
	connections []network.NetworkConnection
}

func GetPeerInstance() *Peer {
	once.Do(func() {
		handlers := map[network.OperationType]messageHandlers.MessageHandler{
			network.JOIN_CHAT:      &messageHandlers.JoinChatHandler{},
			network.SEND_FILE:      &messageHandlers.SendFileHandler{},
			network.SYNC_REQUEST:   &messageHandlers.SyncRequestHandler{},
			network.SYNC_RESPONSE:  &messageHandlers.SyncResponseHandler{},
			network.SET_USERNAME:   &messageHandlers.SetUsernameHandler{},
			network.SEND_MESSAGE:   &messageHandlers.SendMessageHandler{},
			network.CREATE_CHAT:    &messageHandlers.CreateChatHandler{},
			network.INVITE_TO_CHAT: &messageHandlers.InviteToChatHandler{},
			network.LEAVE_CHAT:     &messageHandlers.LeaveChatHandler{},
			network.TEST_MESSAGE:   &messageHandlers.TestMessageHandler{},
			network.TEST_MESSAGE_2: &messageHandlers.TestMessageHandler2{},
		}

		peerInstance = &Peer{
			Chats:       networkChats{},
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

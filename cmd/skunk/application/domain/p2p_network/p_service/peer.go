package p_service

import (
	"errors"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
	"sync"
)

var (
	peerInstance *Peer
	once         sync.Once
)

type Peer struct {
	handlers        map[network.OperationType]messageHandlers.MessageHandler
	connections     []network.NetworkConnection
	securityContext SecurityValidater
	storage         store.NetworkMessageStoragePort
	messageSender   messageHandlers.MessageSender
}

func GetPeerInstance() *Peer {
	once.Do(func() {
		storage := storageSQLiteAdapter.GetInstance("skunk.db")
		securityContext := NewSecurityContext(storage, storage, storage)
		sender := messageHandlers.NewMessageSender(securityContext)

		handlers := map[network.OperationType]messageHandlers.MessageHandler{
			network.JOIN_CHAT:      messageHandlers.NewJoinChatHandler(nil, storage),
			network.SEND_FILE:      messageHandlers.NewSendFileHandler(nil, storage),
			network.SYNC_REQUEST:   messageHandlers.NewSyncRequestHandler(nil, storage),
			network.SYNC_RESPONSE:  messageHandlers.NewSyncResponseHandler(storage),
			network.SET_USERNAME:   messageHandlers.NewSetUsernameHandler(nil, storage),
			network.SEND_MESSAGE:   messageHandlers.NewSendMessageHandler(nil, storage),
			network.INVITE_TO_CHAT: messageHandlers.NewInviteToChatHandler(nil, storage),
			network.LEAVE_CHAT:     messageHandlers.NewLeaveChatHandler(nil, storage),
			network.TEST_MESSAGE:   &messageHandlers.TestMessageHandler{},
			network.TEST_MESSAGE_2: &messageHandlers.TestMessageHandler2{},
		}

		peerInstance = &Peer{
			handlers:        handlers,
			connections:     []network.NetworkConnection{},
			securityContext: securityContext,
			storage:         storage,
			messageSender:   *sender,
		}
	})

	return peerInstance
}

func (p *Peer) AddNetworkConnection(connection network.NetworkConnection) error {
	if len(p.connections) > 0 {
		return errors.New("connection already exists, multiple connections are not supported so far")
	}

	p.connections = append(p.connections, connection)
	connection.SubscribeToNetwork(p)
	p.messageSender.SetNetworkConnection(connection)

	return nil
}

func (p *Peer) RemoveNetworkConnection(connection network.NetworkConnection) {
	for i, c := range p.connections {
		if c == connection {
			p.connections = append(p.connections[:i], p.connections[i+1:]...)
			connection.UnsubscribeFromNetwork(p)
			break
		}
	}
}

func (p *Peer) Notify(message network.Message) error {
	if handler, exists := p.handlers[message.Operation]; exists {
		if !p.securityContext.ValidateIncomingMessage(message) {
			return errors.New("invalid message")
		}

		p.storage.StoreMessage(message)
		return handler.HandleMessage(message)
	}
	return errors.New("invalid message operation")
}

package p_service

import (
	"errors"
	"fmt"
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
	Address         string
	connections     []network.NetworkConnection
	handlers        map[network.OperationType]messageHandlers.MessageHandler
	messageSender   messageHandlers.MessageSender
	securityContext SecurityValidater
	storage         store.NetworkMessageStoragePort
}

func GetPeerInstance() *Peer {
	once.Do(func() {
		storage := storageSQLiteAdapter.GetInstance("skunk.db")
		securityContext := NewSecurityContext(storage, storage, storage)
		sender := messageHandlers.NewMessageSender(securityContext)

		handlers := map[network.OperationType]messageHandlers.MessageHandler{
			network.SEND_MESSAGE:   messageHandlers.NewSendMessageHandler(nil, storage),
			network.SYNC_REQUEST:   messageHandlers.NewSyncRequestHandler(nil, storage),
			network.SYNC_RESPONSE:  messageHandlers.NewSyncResponseHandler(storage),
			network.JOIN_CHAT:      messageHandlers.NewJoinChatHandler(nil, storage),
			network.LEAVE_CHAT:     messageHandlers.NewLeaveChatHandler(nil, storage),
			network.INVITE_TO_CHAT: messageHandlers.NewInviteToChatHandler(nil, storage),
			network.SEND_FILE:      messageHandlers.NewSendFileHandler(nil, storage),
			network.SET_USERNAME:   messageHandlers.NewSetUsernameHandler(nil, storage),
			network.NETWORK_ONLINE: &messageHandlers.NetworkOnlineHandler{},
			network.TEST_MESSAGE:   &messageHandlers.TestMessageHandler{},
			network.TEST_MESSAGE_2: &messageHandlers.TestMessageHandler2{},
		}

		peerInstance = &Peer{
			Address:         "",
			handlers:        handlers,
			connections:     []network.NetworkConnection{},
			securityContext: securityContext,
			storage:         storage,
			messageSender:   *sender,
		}
	})

	return peerInstance
}

// AddNetworkConnection adds a new network connection to the Peer instance
// and subscribes the Peer to the network events.
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
			err := c.UnsubscribeFromNetwork()
			if err != nil {
				fmt.Println(err)
			}
			p.connections = append(p.connections[:i], p.connections[i+1:]...)
			connection.UnsubscribeFromNetwork()
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

		p.storage.StoreMessage(message)
		return handler.HandleMessage(message)
	}
	return errors.New("invalid message operation")
}

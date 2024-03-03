package p_model

import (
	"errors"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkMockAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"sync"
)

var (
	peerInstance *Peer
	once         sync.Once
)

type HandlerMap map[c_model.OperationType]messageHandlers.MessageHandler

type Peer struct {
	Chats      NetworkChats
	handlers   HandlerMap
	connection network.NetworkConnection
}

func GetPeerInstance() *Peer {
	once.Do(func() {
		handlers := map[c_model.OperationType]messageHandlers.MessageHandler{
			c_model.JOIN_CHAT:      &messageHandlers.JoinChatHandler{},
			c_model.SEND_FILE:      &messageHandlers.SendFileHandler{},
			c_model.SYNC_REQUEST:   &messageHandlers.SyncRequestHandler{},
			c_model.SYNC_RESPONSE:  &messageHandlers.SyncResponseHandler{},
			c_model.SET_USERNAME:   &messageHandlers.SetUsernameHandler{},
			c_model.SEND_MESSAGE:   &messageHandlers.SendMessageHandler{},
			c_model.CREATE_CHAT:    &messageHandlers.CreateChatHandler{},
			c_model.INVITE_TO_CHAT: &messageHandlers.InviteToChatHandler{},
			c_model.LEAVE_CHAT:     &messageHandlers.LeaveChatHandler{},
			c_model.TEST_MESSAGE:   &messageHandlers.TestMessageHandler{},
		}

		peerInstance = &Peer{
			Chats:      NetworkChats{},
			handlers:   handlers,
			connection: networkMockAdapter.GetMockConnection(),
		}

		peerInstance.connection.SubscribeToNetwork(peerInstance)
	})

	return peerInstance
}

func (p *Peer) Notify(message c_model.Message) error {
	if handler, exists := p.handlers[message.Operation]; exists {
		return handler.HandleMessage(message) // some form of authentication should be done here
	}
	return errors.New("invalid message operation")
}

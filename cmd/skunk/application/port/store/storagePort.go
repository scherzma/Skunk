package store

import "github.com/scherzma/Skunk/cmd/skunk/application/port/network"

type UserMessageStoragePort interface {
	PeerSetUsername(peerId string, chatId string, username string) error
}

type ChatActionStoragePort interface {
	PeerJoinedChat(timestamp int64, peerId string, chatId string) error
	PeerLeftChat(peerId string, chatId string) error
	ChatCreated(chatName string, chatId string) error // Ensure this line exists
}

type PublicKeyAddress struct {
	Address   string
	PublicKey string
}

type ChatInvitationStoragePort interface {
	InvitedToChat(messageID string, peers []PublicKeyAddress) error
	PeerGotInvitedToChat(peerId string, chatId string) error
	GetInvitations(peerId string) ([]string, error)
}

type SyncStoragePort interface {
	GetMissingInternalMessages(chatId string, inputMessageIDs []string) ([]string, error)
	GetMissingExternalMessages(chatId string, inputMessageIDs []string) ([]string, error)
}

type NetworkMessageStoragePort interface {
	StoreMessage(message network.Message) error
	RetrieveMessage(messageId string) (network.Message, error)
}

type ChatMessage struct {
	Username  string
	Content   string
	MessageId int64
	Timestamp int64
}

type Chat struct {
	ChatId   string
	ChatName string
}

type User struct {
	UserId   string
	Username string
}

type DisplayStoragePort interface {
	GetChats() ([]Chat, error)
	GetUsername(peerID, chatID string) (string, error)
	GetUsersInChat(chatID string) ([]User, error)
	GetPeers() ([]string, error)
	GetChatMessages(chatID string) ([]network.Message, error)
}

type Storage interface {
	UserMessageStoragePort
	ChatActionStoragePort
	ChatInvitationStoragePort
	SyncStoragePort
	NetworkMessageStoragePort
	DisplayStoragePort
}

package store

import "github.com/scherzma/Skunk/cmd/skunk/application/port/network"

type UserMessageStoragePort interface {
	PeerSetUsername(peerId string, chatId string, username string) error
}

type ChatActionStoragePort interface {
	PeerJoinedChat(peerId string, chatId string) error
	PeerLeftChat(peerId string, chatId string, senderId string, date int64) error
	ChatCreated(chatName string, chatId string) error
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
	GetMissingInternalMessages(chatId string, inputMessageIDs []string) []string
	GetMissingExternalMessages(chatId string, inputMessageIDs []string) []string
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
	GetChats() []Chat
	GetUsername(peerId string, chatId string) string
	GetUsersInChat(chatId string) []User
	GetPeers() []string
	GetChatMessages(chatId string) []ChatMessage
}

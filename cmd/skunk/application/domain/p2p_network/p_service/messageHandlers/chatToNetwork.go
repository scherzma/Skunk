package messageHandlers

import (
	"github.com/google/uuid"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
	"time"
)

type ChatToNetwork struct {
	sender  *MessageSender
	storage store.Storage
}

func NewChatToNetwork(sender *MessageSender, chatActionStorage store.Storage) *ChatToNetwork {
	return &ChatToNetwork{
		sender:  sender,
		storage: chatActionStorage,
	}
}

func (c *ChatToNetwork) CreateChat(chatId string, chatName string) error {
	err := c.storage.ChatCreated(chatName, chatId)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatToNetwork) JoinChat(chatId string) error {
	timestamp := time.Now().UnixNano()

	err := c.storage.PeerJoinedChat(timestamp, "self", chatId)
	if err != nil {
		return err
	}

	message := network.Message{
		Id:              uuid.New().String(),
		Timestamp:       timestamp,
		Content:         "",
		SenderID:        "self",
		ReceiverID:      "?",
		SenderAddress:   "self",
		ReceiverAddress: "?",
		ChatID:          chatId,
		Operation:       network.JOIN_CHAT,
	}

	c.storage.StoreMessage(message)
	c.sender.SendMessage(message)
	return nil
}

func (c *ChatToNetwork) LeaveChat(chatId string) error {
	err := c.storage.PeerLeftChat("self", chatId)
	if err != nil {
		return err
	}

	message := network.Message{
		Id:              uuid.New().String(),
		Timestamp:       time.Now().UnixNano(),
		Content:         "",
		SenderID:        "self",
		ReceiverID:      "?",
		SenderAddress:   "self",
		ReceiverAddress: "?",
		ChatID:          chatId,
		Operation:       network.LEAVE_CHAT,
	}

	c.storage.StoreMessage(message)
	c.sender.SendMessage(message)
	return nil
}

func (c *ChatToNetwork) InviteToChat(chatId string, peerId string) error {

	err := c.storage.PeerGotInvitedToChat(peerId, chatId)
	if err != nil {
		return err
	}

	message := network.Message{
		Id:              uuid.New().String(),
		Timestamp:       time.Now().UnixNano(),
		Content:         "", //TODO fill
		SenderID:        "self",
		ReceiverID:      "?",
		SenderAddress:   "self",
		ReceiverAddress: peerId,
		ChatID:          chatId,
		Operation:       network.INVITE_TO_CHAT,
	}

	c.storage.StoreMessage(message)
	c.sender.SendMessage(message)

	return nil
}

func (c *ChatToNetwork) SendFileToChat(chatId string, filePath string) error {

	message := network.Message{
		Id:              uuid.New().String(),
		Timestamp:       time.Now().UnixNano(),
		Content:         "", //TODO fill
		SenderID:        "self",
		ReceiverID:      "?",
		SenderAddress:   "self",
		ReceiverAddress: "?",
		ChatID:          chatId,
		Operation:       network.INVITE_TO_CHAT,
	}

	c.storage.StoreMessage(message)
	c.sender.SendMessage(message)
	return nil
}

func (c *ChatToNetwork) SetUsernameInChat(chatId string, username string) error {
	return nil
}

func (c *ChatToNetwork) SendMessageToChat(chatId string, message string) error {
	return nil
}

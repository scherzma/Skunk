package test

import "fmt"

type MockChatLogic struct {
	LastSenderId    string
	LastChatId      string
	LastChatName    string
	LastChatMembers []string
}

func (m *MockChatLogic) ReceiveMessage(senderId string, chatId string, message string) error {
	return nil
}

func (m *MockChatLogic) ReceiveChatInvitation(senderId string, chatId string, chatName string, chatMembers []string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.LastChatName = chatName
	m.LastChatMembers = chatMembers
	fmt.Printf("Invitation received from %s to join chat %s (%s) with members %v\n", senderId, chatId, chatName, chatMembers)
	return nil
}

func (m *MockChatLogic) PeerLeavesChat(senderId string, chatId string) error {
	return nil
}

func (m *MockChatLogic) PeerJoinsChat(senderId string, chatId string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	fmt.Printf("Peer %s joined chat %s\n", senderId, chatId)
	return nil
}
func (m *MockChatLogic) ReceiveFile(senderId string, chatId string, filePath string) error {
	return nil
}

func (m *MockChatLogic) PeerSetsUsername(senderId string, chatId string, username string) error {
	return nil
}

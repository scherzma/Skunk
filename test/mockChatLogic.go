package test

import "fmt"

type MockChatLogic struct {
	LastSenderId    string
	LastChatId      string
	LastChatName    string
	LastChatMembers []string
	LastFileName    string
	LastFileSize    int
	LastFileData    string
	LastMessage     string
	LastUsername    string
}

func (m *MockChatLogic) ReceiveMessage(senderId string, chatId string, message string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.LastMessage = message
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
	m.LastSenderId = senderId
	m.LastChatId = chatId
	return nil
}

func (m *MockChatLogic) PeerJoinsChat(senderId string, chatId string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	fmt.Printf("Peer %s joined chat %s\n", senderId, chatId)
	return nil
}
func (m *MockChatLogic) ReceiveFile(senderId string, chatId string, filePath string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.LastFileName = filePath
	fmt.Printf("Received file from %s in chat %s: %s\n", senderId, chatId, filePath)
	return nil
}

func (m *MockChatLogic) PeerSetsUsername(senderId string, chatId string, username string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.LastUsername = username
	fmt.Printf("Peer %s set username to %s in chat %s\n", senderId, username, chatId) // Add a print statement for debugging
	return nil
}

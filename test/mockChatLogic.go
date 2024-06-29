package test

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
	LogEntries      []string
}

func (m *MockChatLogic) log(message string) {
	m.LogEntries = append(m.LogEntries, message)
}

func (m *MockChatLogic) ReceiveMessage(senderId string, chatId string, message string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.LastMessage = message
	m.log("ReceiveMessage called")
	return nil
}

func (m *MockChatLogic) ReceiveChatInvitation(senderId string, chatId string, chatName string, chatMembers []string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.LastChatName = chatName
	m.LastChatMembers = chatMembers
	m.log("ReceiveChatInvitation called")
	return nil
}

func (m *MockChatLogic) PeerLeavesChat(senderId string, chatId string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.log("PeerLeavesChat called")
	return nil
}

func (m *MockChatLogic) PeerJoinsChat(senderId string, chatId string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.log("PeerJoinsChat called")
	return nil
}

func (m *MockChatLogic) ReceiveFile(senderId string, chatId string, filePath string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.LastFileName = filePath
	m.log("ReceiveFile called")
	return nil
}

func (m *MockChatLogic) PeerSetsUsername(senderId string, chatId string, username string) error {
	m.LastSenderId = senderId
	m.LastChatId = chatId
	m.LastUsername = username
	m.log("PeerSetsUsername called")
	return nil
}

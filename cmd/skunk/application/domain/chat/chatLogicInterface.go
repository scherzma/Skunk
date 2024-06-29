package chat

type ChatLogic interface {
	ReceiveMessage(senderId string, chatId string, message string) error
	ReceiveChatInvitation(senderId string, chatId string, chatName string, chatMembers []string) error
	PeerLeavesChat(senderId string, chatId string) error
	PeerJoinsChat(senderId string, chatId string) error
	ReceiveFile(senderId string, chatId string, filePath string) error
	PeerSetsUsername(senderId string, chatId string, username string) error
}

package p2p_network

type NetworkLogic interface {
	// TODO: Implement
	CreateChat(chatId string, chatName string) error
	JoinChat(chatId string) error
	LeaveChat(chatId string) error
	InviteToChat(chatId string, peerId string) error
	SendFileToChat(chatId string, filePath string) error
	SetUsernameInChat(chatId string, username string) error
	SendMessageToChat(chatId string, message string) error
}

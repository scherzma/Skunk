package p_model

type NetworkChats struct {
	chatMap map[string]NetworkChatMessages
}

func (n *NetworkChats) AddChat(chatId string) {
	n.chatMap[chatId] = *NewNetworkChatMessages()
}

func (n *NetworkChats) GetChat(chatId string) NetworkChatMessages {
	return n.chatMap[chatId]
}

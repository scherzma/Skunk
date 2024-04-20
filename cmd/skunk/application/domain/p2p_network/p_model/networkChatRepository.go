package p_model

type networkChats struct {
	chatMap map[string]NetworkChatMessages
}

var instance *networkChats

func GetNetworkChatsInstance() *networkChats {
	if instance == nil {
		instance = &networkChats{
			chatMap: make(map[string]NetworkChatMessages),
		}
	}
	return instance
}

func (n *networkChats) AddChat(chatId string) {
	n.chatMap[chatId] = *NewNetworkChatMessages()
}

func (n *networkChats) GetChat(chatId string) NetworkChatMessages {
	return n.chatMap[chatId]
}

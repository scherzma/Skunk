package p_model

type NetworkChats struct {
	chatMap map[string]NetworkChatMessages
}

var instance *NetworkChats

func GetNetworkChatsInstance() *NetworkChats {
	if instance == nil {
		instance = &NetworkChats{
			chatMap: make(map[string]NetworkChatMessages),
		}
	}
	return instance
}

func (n *NetworkChats) AddChat(chatId string) {
	n.chatMap[chatId] = *NewNetworkChatMessages()
}

func (n *NetworkChats) GetChat(chatId string) NetworkChatMessages {
	if _, exists := n.chatMap[chatId]; !exists { //TODO change default behavior; probably should not create a new chat if it does not exist
		n.AddChat(chatId)
	}
	return n.chatMap[chatId]
}

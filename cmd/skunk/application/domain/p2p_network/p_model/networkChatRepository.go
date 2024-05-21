package p_model

// NetworkChats represents a collection of chat messages organized by chat ID.
// It is implemented as a singleton instance.
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

// AddChat adds a new chat with the specified chat ID to the NetworkChats.
// If the chat already exists, it will be overwritten.
func (n *NetworkChats) AddChat(chatId string) {
	if _, ok := n.chatMap[chatId]; ok {
		return
	}
	n.chatMap[chatId] = *NewNetworkChatMessages()
}

// GetChat returns the NetworkChatMessages for the specified chat ID.
// If the chat doesn't exist, it creates a new chat with the given ID.
// TODO: Change the default behavior; probably should not create a new chat if it doesn't exist.
func (n *NetworkChats) GetChat(chatId string) NetworkChatMessages {
	if _, exists := n.chatMap[chatId]; !exists { //TODO change default behavior; probably should not create a new chat if it does not exist
		n.AddChat(chatId)
	}
	return n.chatMap[chatId]
}

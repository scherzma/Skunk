package chat

type ChatLogic interface {
	RecieveMessage() error
	ReceiveChatInvitation() error
	PeerLeavesChat() error
	PeerJoinsChat() error
	ReceiveFile() error
	PeerSetsUsername() error
}

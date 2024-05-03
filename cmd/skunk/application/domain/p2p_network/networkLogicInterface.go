package p2p_network

type NetworkLogic interface {
	// TODO: Implement
	CreateChat() error
	JoinChat() error
	LeaveChat() error
	InviteToChat() error
	SendFileToChat() error
	SetUsername() error
}

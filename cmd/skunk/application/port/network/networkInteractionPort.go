package network

type OperationType int

const (
	SEND_MESSAGE   OperationType = iota
	SYNC_REQUEST   OperationType = iota
	SYNC_RESPONSE  OperationType = iota
	CREATE_CHAT    OperationType = iota
	JOIN_CHAT      OperationType = iota
	LEAVE_CHAT     OperationType = iota
	INVITE_TO_CHAT OperationType = iota
	SEND_FILE      OperationType = iota
	SET_USERNAME   OperationType = iota
	TEST_MESSAGE   OperationType = iota
	TEST_MESSAGE_2 OperationType = iota
)

type Message struct {
	Id        string
	Timestamp int64
	Content   string
	FromUser  string // UserID
	ChatID    string // ChatID
	Operation OperationType
}

type NetworkObserver interface {
	Notify(message Message) error
}

type NetworkConnection interface {
	SubscribeToNetwork(observer NetworkObserver) error
	UnsubscribeFromNetwork(observer NetworkObserver) error
	SendMessageToNetworkPeer(address string, message Message) error
}

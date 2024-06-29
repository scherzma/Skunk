package network

// OperationType represents the type of operation performed in a network message.
type OperationType int

const (
	SEND_MESSAGE   OperationType = iota
	SYNC_REQUEST   OperationType = iota
	SYNC_RESPONSE  OperationType = iota
	JOIN_CHAT      OperationType = iota
	LEAVE_CHAT     OperationType = iota
	INVITE_TO_CHAT OperationType = iota
	SEND_FILE      OperationType = iota
	SET_USERNAME   OperationType = iota
	USER_OFFLINE   OperationType = iota
	NETWORK_ONLINE OperationType = iota
	TEST_MESSAGE   OperationType = iota
	TEST_MESSAGE_2 OperationType = iota
)

// Message represents a network message exchanged between peers.
type Message struct {
	Id              string
	Timestamp       int64
	Content         string
	SenderID        string
	ReceiverID      string
	SenderAddress   string
	ReceiverAddress string
	ChatID          string
	Operation       OperationType
}

// NetworkObserver is an interface that defines the contract for observing network events.
// Types that implement this interface can be notified of incoming network messages.
type NetworkObserver interface {
	Notify(message Message) error
}

// NetworkConnection is an interface that defines the contract for a network connection.
// It provides methods for subscribing/unsubscribing observers and sending messages to network peers.
type NetworkConnection interface {
	SubscribeToNetwork(observer NetworkObserver) error
	UnsubscribeFromNetwork() error
	SendMessageToNetworkPeer(message Message) error
}

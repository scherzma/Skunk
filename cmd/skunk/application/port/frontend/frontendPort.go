package frontend

// FrontendMessage represents a message sent from the frontend to the backend
type FrontendMessage struct {
	Timestamp int64
	Content   string
	FromUser  string // UserID
	ChatID    string // ChatID
	Operation OperationType
}

// OperationType represents the different types of operations that can be performed
type OperationType int

const (
	SEND_MESSAGE   OperationType = iota
	CREATE_CHAT    OperationType = iota
	JOIN_CHAT      OperationType = iota
	LEAVE_CHAT     OperationType = iota
	INVITE_TO_CHAT OperationType = iota
	SEND_FILE      OperationType = iota
	SET_USERNAME   OperationType = iota
	TEST_MESSAGE   OperationType = iota
)

// FrontendObserver is an interface for observing messages from the frontend
type FrontendObserver interface {
	Notify(message FrontendMessage) error
}

// Frontend is an interface for interacting with the frontend
type Frontend interface {
	SubscribeToFrontend(observer FrontendObserver) error
	UnsubscribeFromFrontend(observer FrontendObserver) error
	SendToFrontend(message FrontendMessage) error // TODO: change to FrontendMessage
}

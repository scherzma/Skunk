package frontend

type FrontendMessage struct {
	Timestamp int64
	Content   string
	FromUser  string // UserID
	ChatID    string // ChatID
	Operation OperationType
}

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

type FrontendObserver interface {
	Notify(message FrontendMessage) error
}

type Frontend interface {
	SubscribeToFrontend(observer FrontendObserver) error
	UnsubscribeFromFrontend(observer FrontendObserver) error
	SendToFrontend(message FrontendMessage) error // TODO: change to FrontendMessage
}

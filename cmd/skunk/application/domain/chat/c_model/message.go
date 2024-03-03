package c_model

type OperationType int

const (
	SEND_MESSAGE    OperationType = iota
	SYNC_REQUEST    OperationType = iota
	SYNC_RESPONSE   OperationType = iota
	CREATE_GROUP    OperationType = iota
	JOIN_GROUP      OperationType = iota
	LEAVE_GROUP     OperationType = iota
	INVITE_TO_GROUP OperationType = iota
	SEND_FILE       OperationType = iota
	SET_USERNAME    OperationType = iota
	TEST_MESSAGE    OperationType = iota
)

type Message struct {
	Id        string
	Timestamp int64
	Content   string
	From      User // should probably use User ID instead of a direct reference
	To        User // shoud honestly be anything but a user... like a chat...
	Operation OperationType
}

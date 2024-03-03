package messageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"

type SyncRequestHandler struct {
}

func (s *SyncRequestHandler) HandleMessage(message c_model.Message) error {
	//TODO implement
	return nil
}

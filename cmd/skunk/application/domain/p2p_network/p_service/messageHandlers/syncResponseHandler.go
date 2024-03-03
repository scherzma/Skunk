package messageHandlers

import "github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"

type SyncResponseHandler struct{}

func (s *SyncResponseHandler) HandleMessage(message c_model.Message) error {
	//TODO implement
	return nil
}

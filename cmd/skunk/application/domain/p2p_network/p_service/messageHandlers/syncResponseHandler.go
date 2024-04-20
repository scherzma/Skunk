package messageHandlers

import (
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type SyncResponseHandler struct{}

func (s *SyncResponseHandler) HandleMessage(message network.Message) error {
	//TODO implement
	return nil
}

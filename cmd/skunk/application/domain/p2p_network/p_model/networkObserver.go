package p_model

import "github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"

type NetworkObserver struct {
}

func (n *NetworkObserver) Notify(message c_model.Message) error {
	//TODO implement
	return nil
}

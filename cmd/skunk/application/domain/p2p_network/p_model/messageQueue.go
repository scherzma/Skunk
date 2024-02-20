package p_model

import "github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"

type MessageQueue struct {
	messageQueue []c_model.Message // shouldn't depend on the message model?
}

func sendMessages() {

}

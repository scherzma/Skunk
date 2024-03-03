package p_model

import (
	"errors"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/chat/c_model"
)

type MessageHandler interface {
	HandleMessage(message c_model.Message) error
}

type HandlerMap map[c_model.OperationType]MessageHandler

type Peer struct {
	Messages *treemap.Map
	Handlers HandlerMap
}

func NewPeer() *Peer {
	return &Peer{
		Messages: treemap.NewWith(MessageComparator),
		Handlers: make(HandlerMap),
	}
}

func MessageComparator(a, b interface{}) int {
	aM, ok := a.(c_model.Message)
	if !ok {
		return 0
	}
	bM, ok := b.(c_model.Message)
	if !ok {
		return 0
	}
	return int(aM.Timestamp - bM.Timestamp)
}

func (p *Peer) RecieveMessage(message c_model.Message) error {
	if handler, ok := p.Handlers[message.Operation]; ok {
		return handler.HandleMessage(message)
	}
	return errors.New("Invalid message operation")
}

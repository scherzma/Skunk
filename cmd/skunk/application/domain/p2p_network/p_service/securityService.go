package p_service

import "github.com/scherzma/Skunk/cmd/skunk/application/port/network"

type SecurityContext struct {
}

// SecurityContext is a service that provides security checks for the network
// It should be possible to implement all security checks in this service

func (s *SecurityContext) ValidateOutgoingMessage(message network.Message) bool {
	return true
}

func (s *SecurityContext) ValidateIncomingMessage(message network.Message) bool {
	return true
}

func (s *SecurityContext) ValidatePeer(peer string) bool {
	return true
}

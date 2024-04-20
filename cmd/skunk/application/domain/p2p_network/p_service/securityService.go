package p_service

import "github.com/scherzma/Skunk/cmd/skunk/application/port/network"

func ValidateMessage(message network.Message) bool {
	return true
}

func ValidatePeer(peer string) bool {
	return true
}

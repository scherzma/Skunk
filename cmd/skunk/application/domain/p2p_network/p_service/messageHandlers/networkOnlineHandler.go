package messageHandlers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type NetworkOnlineHandler struct{}

func (s *NetworkOnlineHandler) HandleMessage(message network.Message) error {
	peer := GetPeerInstance()

	var peerAddress string
	err := json.Unmarshal([]byte(message.Content), &peerAddress)
	if err != nil {
		return fmt.Errorf("error unmarshalling peer address: %v", err)
	}

	if !strings.HasPrefix(peerAddress, "ws://") || !strings.Contains(peerAddress, ".onion:") {
		return fmt.Errorf("network has sent incorrectly formatted peer address")
	}
	peer.Address = peerAddress
	return nil
}

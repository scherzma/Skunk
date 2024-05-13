package messageHandlers

import (
    "fmt"
    "strings"

    "github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

type NetworkOnlineHandler struct {}

func (s *NetworkOnlineHandler) HandleMessage(message network.Message) error {
    peer := GetPeerInstance()

    peerAddress := message.Content
    if !strings.HasPrefix(peerAddress, "ws://") || !strings.Contains(peerAddress, ".onion:") {
        return fmt.Errorf("network has sent incorrectly formatted peer address")
    }
    peer.Address = peerAddress
    return nil
}

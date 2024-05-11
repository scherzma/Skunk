package test

import (
    "testing"
    "time"

    "github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkAdapter"
    "github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_model"
    "github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
    "github.com/scherzma/Skunk/cmd/skunk/application/port/network"
    "github.com/scherzma/Skunk/cmd/skunk/util/timestamp"
    "github.com/scherzma/Skunk/cmd/skunk/util/uuid"
)

func TestNetworkAdapter(t *Testing.T) {
    testMessage := network.Message {
        Id: uuid.UUID(),
        Timestamp: timestamp.CurrentTimeMillis(),
        Content: "Hello World!",
        FromUser: "Alice",
        ChatID: "1",
        Operation: network.TEST_MESSAGE,
    }

    peerInstance := messageHandlers.GetPeerInstance()

    networkConnection := networkAdapter.NetworkConnection()
    peerInstance.AddNetworkConnection(networkConnection)

    peer.SendMessageToNetworkPeer(, testMessage)

    peer.RemoveNetworkConnection(networkConnection)
}

package test

import (
	"testing"
    "time"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/networkAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/peer"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/tor"
	"github.com/scherzma/Skunk/cmd/skunk/application/domain/p2p_network/p_service/messageHandlers"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/util"
	"github.com/stretchr/testify/assert"
)

func TestNetworkAdapter(t *testing.T) {
	testMessage := network.Message{
		Id:        util.UUID(),
		Timestamp: util.CurrentTimeMillis(),
		Content:   "",
		FromUser:  "Alice",
		ChatID:    "1",
		Operation: network.TEST_MESSAGE,
	}

	peerInstance := messageHandlers.GetPeerInstance()
	networkConnection := networkAdapter.NewAdapter()
	peerInstance.AddNetworkConnection(networkConnection)

    conf := &tor.TorConfig{
        DataDir: "data-dir-local",
        SocksPort: "9052",
        LocalPort: "1110",
        RemotePort: "2220",
        DeleteDataDirOnClose: false,
        UseEmbedded: false,
    }

    torInstance, _ := tor.NewTor(conf)
    torInstance.StartTor()
    onion, _ := torInstance.StartHiddenService()

    time.Sleep(10 * time.Second)

    peerNetworkInstance, _ := peer.NewPeer(onion.ID+".onion", "1110", "2220", "127.0.0.1:9052")
    defer peerNetworkInstance.Shutdown()

    peerNetworkInstance.Listen(onion)
    time.Sleep(1 * time.Second)

    messageCh := make(chan string)
    errorCh := make(chan error)

    go peerNetworkInstance.ReadMessages(messageCh, errorCh)

	err := peerInstance.SendMessageToNetworkPeer(peerNetworkInstance.Address, testMessage)
	assert.NoError(t, err)

    time.Sleep(10 * time.Second)

    select {
        case msg := <-messageCh:
            t.Log(msg)
        case err := <-errorCh:
            assert.NoError(t, err)
        default:
    }

    torInstance.StopTor()
	peerInstance.RemoveNetworkConnection(networkConnection)
}


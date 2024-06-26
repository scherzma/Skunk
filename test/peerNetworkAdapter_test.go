package test

import (
	"encoding/json"
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

/*
We have to test sending and receiving a message separately, because the application peer is a singleton and therefore can only exist once and because we cannot send messages from the same socket to the same socket, we need an application peer once and the network peer once, then the network peer acts once as a client and once as a server.
Since the embedded tor process is already used in the application peer, the network peer must run over the locally installed Tor version.
*/

func TestNetworkAdapterSendMessage(t *testing.T) {
	testMessage := network.Message{
		Id:        util.UUID(),
		Timestamp: util.CurrentTimeMillis(),
		Content:   "",
		FromUser:  "Alice",
		ChatID:    "1",
		Operation: network.TEST_MESSAGE,
	}
	testMessageJson, _ := json.Marshal(testMessage)

	peerInstance := messageHandlers.GetPeerInstance()
	networkConnection := networkAdapter.NewAdapter()
	peerInstance.AddNetworkConnection(networkConnection)

	conf := &tor.TorConfig{
		DataDir:              "data-dir-local",
		SocksPort:            "9052",
		LocalPort:            "1110",
		RemotePort:           "2220",
		DeleteDataDirOnClose: false,
		UseEmbedded:          false,
	}

	torInstance, _ := tor.NewTor(conf)
	torInstance.StartTor()
	onion, _ := torInstance.StartHiddenService()

	peerNetworkInstance, _ := peer.NewPeer(onion.ID+".onion", "1110", "2220", "127.0.0.1:9052")
	peerNetworkInstance.Listen(onion)

	defer func() {
		// sleep because shutdown takes one second
		peerNetworkInstance.Shutdown()
		time.Sleep(2 * time.Second)
		torInstance.StopTor()

		peerInstance.RemoveNetworkConnection(networkConnection)
	}()

	messageCh := make(chan string)
	errorCh := make(chan error)

	go peerNetworkInstance.ReadMessages(messageCh, errorCh)

	err := peerInstance.SendMessageToNetworkPeer(peerNetworkInstance.Address, testMessage)
	assert.NoError(t, err)

	time.Sleep(1 * time.Minute)

	select {
	case msg := <-messageCh:
		t.Log(msg)
		assert.Equal(t, string(testMessageJson), msg)
	case err := <-errorCh:
		assert.NoError(t, err)
	default:
		assert.Falsef(t, true, "message could not be received")
	}
}

func TestNetworkAdapterReceiveMessage(t *testing.T) {
	// We don't have an exact insight from this code into what is happening internally within the application peer, so we use a little trick by sending a NETWORK_ONLINE message which, if everything works, triggers the application peer to change its address, which we have access to in this code.
	testMessage := network.Message{
		Id:        util.UUID(),
		Timestamp: util.CurrentTimeMillis(),
		Content:   `"ws://testworked.onion:1111"`,
		FromUser:  "Bob",
		Operation: network.NETWORK_ONLINE,
	}
	testMessageJson, _ := json.Marshal(testMessage)

	peerInstance := messageHandlers.GetPeerInstance()
	networkConnection := networkAdapter.NewAdapter()
	peerInstance.AddNetworkConnection(networkConnection)

	conf := &tor.TorConfig{
		DataDir:              "data-dir-local-send",
		SocksPort:            "9053",
		LocalPort:            "3330",
		RemotePort:           "4440",
		DeleteDataDirOnClose: false,
		UseEmbedded:          false,
	}

	torInstance, _ := tor.NewTor(conf)
	torInstance.StartTor()
	onion, _ := torInstance.StartHiddenService()

	peerNetworkInstance, _ := peer.NewPeer(onion.ID+".onion", "3330", "4440", "127.0.0.1:9053")
	peerNetworkInstance.Listen(onion)

	defer func() {
		peerNetworkInstance.Shutdown()
		// sleep because shutdown takes one second
		time.Sleep(2 * time.Second)
		torInstance.StopTor()

		peerInstance.RemoveNetworkConnection(networkConnection)
	}()

	err := peerNetworkInstance.Connect(peerInstance.Address)
	assert.NoError(t, err)

	err = peerNetworkInstance.SetWriteConn(peerInstance.Address)
	assert.NoError(t, err)

	err = peerNetworkInstance.WriteMessage(string(testMessageJson))
	assert.NoError(t, err)

	time.Sleep(1 * time.Minute)

	assert.Equal(t, "ws://testworked.onion:1111", peerInstance.Address)
}

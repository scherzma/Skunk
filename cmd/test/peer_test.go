package test

import (
	"fmt"
	"testing"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/peer"

	"github.com/stretchr/testify/assert"
)

func TestPeerSendMessageWorkflow(t *testing.T) {
	peer1, err := peer.NewPeer("127.0.0.1", "1111", "")
	if err != nil {
		t.Errorf("Error creating peer one %v", err)
	}

	assert.Equal(t, peer1.Hostname, "127.0.0.1")
	assert.Equal(t, peer1.Port, "1111")
	assert.Equal(t, peer1.ProxyAddr, "")

	peer2, err := peer.NewPeer("127.0.0.1", "6969", "")
	if err != nil {
		t.Errorf("Error creating peer one %v", err)
	}

	assert.Equal(t, peer2.Hostname, "127.0.0.1")
	assert.Equal(t, peer2.Port, "6969")
	assert.Equal(t, peer2.ProxyAddr, "")

	t.Log("Created both peers!")

	peer1.Listen()
	peer2.Listen()

	t.Log("Both peers are listening")

	err_connect := peer1.Connect(fmt.Sprintf("ws://%s:%s", peer2.Hostname, peer2.Port))
	if err_connect != nil {
		t.Errorf("Error connecting to peer two %v", err_connect)
	}

	err_message := peer1.WriteMessage("Hello Peer Two!")
	if err_message != nil {
		t.Errorf("Error sending message to peer two %v", err)
	}

	v, err := peer2.ReadMessage()
	if err != nil {
		t.Errorf("Error reading message peer two %v", err)
	}

	t.Log("Peer One -> ", v)
	assert.Equal(t, v.(string), "Hello Peer Two!")

	peer1.Shutdown()
	peer2.Shutdown()

	t.Log("Both peers are shut down")
}

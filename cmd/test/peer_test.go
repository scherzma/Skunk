package test

import (
    "time"
	"testing"
    "context"

	"github.com/stretchr/testify/assert"

    "nhooyr.io/websocket"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/peer"
)


func TestNewPeer(t *testing.T) {
    // In the future this should include a test with a proxy

    t.Run("initialize with valid parameters", func(t *testing.T) {
        peerInstance, err := peer.NewPeer("127.0.0.1", "8080", "")
        assert.NoError(t, err)
        assert.NotNil(t, peerInstance)
        assert.Equal(t, peerInstance.Hostname, "127.0.0.1")
        assert.Equal(t, peerInstance.Port, "8080")
        assert.Equal(t, peerInstance.ProxyAddr, "")
    })
}

func TestListen(t *testing.T) {
    peerInstance, err := peer.NewPeer("127.0.0.1", "8080", "")
    assert.NoError(t, err)

    peerInstance.Listen()
    time.Sleep(1 * time.Second)

    conn, _, err := websocket.Dial(context.Background(), "ws://127.0.0.1:8080", nil)
    assert.NoError(t, err)
    assert.NotNil(t, conn)

    peerInstance.Shutdown()
    defer conn.Close(websocket.StatusNormalClosure, "test completed")
}

func TestPeerSetWriteConn(t *testing.T) {
    peer1, err := peer.NewPeer("127.0.0.1", "1111", "")
    assert.NoError(t, err)

    peer2, err := peer.NewPeer("127.0.0.1", "10000", "")
    assert.NoError(t, err)

    peer1.Listen()
    time.Sleep(1 * time.Second)

    err = peer2.Connect(peer1.Address)
    assert.NoError(t, err)
    time.Sleep(1 * time.Second)

    err = peer2.SetWriteConn(peer1.Address)
    assert.NoError(t, err)

    peer1.Shutdown()
    peer2.Shutdown()
}

func TestPeerSendMessage(t *testing.T) {
	peer1, err := peer.NewPeer("127.0.0.1", "2222", "")
    assert.NoError(t, err)

	peer2, err := peer.NewPeer("127.0.0.1", "6969", "")
    assert.NoError(t, err)

	peer1.Listen()
    time.Sleep(1 * time.Second)
	peer2.Listen()
    time.Sleep(1 * time.Second)

	err = peer1.Connect(peer2.Address)
    assert.NoError(t, err)
    time.Sleep(1 * time.Second)

    err = peer1.SetWriteConn(peer2.Address)
    assert.NoError(t, err)

    messageCh := make(chan string)
    errorCh := make(chan error)
	go peer2.ReadMessages(messageCh, errorCh)

	err = peer1.WriteMessage("Hello Peer Two!")
    assert.NoError(t, err)

    select {
    case msg := <-messageCh:
        assert.Equal(t, msg, "From ws://127.0.0.1:2222: Hello Peer Two!")
    case err := <-errorCh:
        assert.NoError(t, err)
    }

	err = peer1.WriteMessage("Hello again Peer Two!")
    assert.NoError(t, err)

    select {
    case msg := <-messageCh:
        assert.Equal(t, msg, "From ws://127.0.0.1:2222: Hello again Peer Two!")
    case err := <-errorCh:
        assert.NoError(t, err)
    }

	peer1.Shutdown()
	peer2.Shutdown()
}


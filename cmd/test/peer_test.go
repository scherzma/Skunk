package test

import (
    "time"
	"testing"
    "context"

	"github.com/stretchr/testify/assert"
    "nhooyr.io/websocket"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/peer"
)

// In the future this should include tests with a proxy

func TestNewPeer(t *testing.T) {
    peerInstance, err := peer.NewPeer("127.0.0.1", "8080", "")
    assert.NoError(t, err)
    assert.NotNil(t, peerInstance)
    assert.Equal(t, peerInstance.Hostname, "127.0.0.1")
    assert.Equal(t, peerInstance.Port, "8080")
    assert.Equal(t, peerInstance.ProxyAddr, "")
}

func TestListen(t *testing.T) {
    peerInstance, err := peer.NewPeer("127.0.0.1", "8080", "")
    assert.NoError(t, err)

    peerInstance.Listen()
    time.Sleep(1 * time.Second)

    // Connecting to peer address should work
    conn, _, err := websocket.Dial(context.Background(), "ws://127.0.0.1:8080", nil)
    assert.NoError(t, err)
    assert.NotNil(t, conn)

    peerInstance.Shutdown()
    defer conn.Close(websocket.StatusNormalClosure, "test completed")
}

func TestConnect(t *testing.T) {
    peer1, err := peer.NewPeer("127.0.0.1", "1234", "")
    peer2, err := peer.NewPeer("127.0.0.1", "4321", "")

    peer1.Listen()
    peer2.Listen()
    time.Sleep(1 * time.Second)

    // First connect from peer1 to peer2.
    err = peer1.Connect(peer2.Address)
    assert.NoError(t, err)

    // Second connect from peer1 to peer2 should return an error because they have already been connected.
    err = peer1.Connect(peer2.Address)
    assert.Error(t, err)

    // Connect from peer1 to peer1 should not work.
    err = peer1.Connect(peer1.Address)
    assert.Error(t, err)

    // Connect from peer2 to peer1 should return an error because they have already been connected.
    err = peer2.Connect(peer1.Address)
    assert.Error(t, err)

    // Should return an error, because "" is not a valid address
    err = peer2.Connect("")
    assert.Error(t, err)

    peer1.Shutdown()
    peer2.Shutdown()
}

func TestPeerSetWriteConn(t *testing.T) {
    peer1, err := peer.NewPeer("127.0.0.1", "1111", "")
    peer2, err := peer.NewPeer("127.0.0.1", "10000", "")

    peer1.Listen()
    time.Sleep(1 * time.Second)

    err = peer2.Connect(peer1.Address)
    time.Sleep(1 * time.Second)

    // First time setting the write conn to peer1.Address should work
    err = peer2.SetWriteConn(peer1.Address)
    assert.NoError(t, err)

    // Setting the write conn to ones own address should not work
    err = peer2.SetWriteConn(peer2.Address)
    assert.Error(t, err)

    // Setting the write conn to an address the peer is not connected to should not work
    err = peer2.SetWriteConn("ws://127.0.0.1:9999")
    assert.Error(t, err)

    // Setting the write conn to "" should not work
    err = peer2.SetWriteConn("")
    assert.Error(t, err)

    peer1.Shutdown()
    peer2.Shutdown()
}

func TestPeerReadMessages(t * testing.T) {
    peer1, _ := peer.NewPeer("127.0.0.1", "2222", "")
    peer2, _ := peer.NewPeer("127.0.0.1", "3333", "")
    peer3, _ := peer.NewPeer("127.0.0.1", "4444", "")
    peer4, _ := peer.NewPeer("127.0.0.1", "5555", "")
    peer5, _ := peer.NewPeer("127.0.0.1", "6666", "")

    peer1.Listen()
    time.Sleep(1 * time.Second)

    address := peer1.Address
    peer2.Connect(address)
    peer3.Connect(address)
    peer4.Connect(address)
    peer5.Connect(address)

    messageCh := make(chan string)
    errorCh := make(chan error)
    go peer1.ReadMessages(messageCh, errorCh)

    peer2.SetWriteConn(address)
    peer3.SetWriteConn(address)
    peer4.SetWriteConn(address)
    peer5.SetWriteConn(address)

    peer2.WriteMessage("Hello World!")
    peer3.WriteMessage("This is the story of my life")
    peer4.WriteMessage("Just do it!")
    peer5.WriteMessage("ABCDEFGHIJKLM")

    peer2.WriteMessage("Hello Proxima Centauri!")
    peer3.WriteMessage("Are you alright?")
    peer4.WriteMessage("No I'm all left!")
    peer5.WriteMessage("XYZ")

    peer2.WriteMessage("Recursive ...")
    peer2.WriteMessage("Recursive ...")
    peer2.WriteMessage("Recursive ...")
    peer2.WriteMessage("Recursive ...")

    // wait until all messages have been sent and received
    time.Sleep(10 * time.Second)

    // check that no error occured in this time
    select {
    case err := <-errorCh:
        assert.NoError(t, err)
    default:
    }

    peer1.Shutdown()
    peer2.Shutdown()
    peer3.Shutdown()
    peer4.Shutdown()
    peer5.Shutdown()
}

func TestPeerSendMessage(t *testing.T) {
	peer1, err := peer.NewPeer("127.0.0.1", "2222", "")
	peer2, err := peer.NewPeer("127.0.0.1", "6969", "")

	peer1.Listen()
    time.Sleep(1 * time.Second)
	peer2.Listen()
    time.Sleep(1 * time.Second)

	err = peer1.Connect(peer2.Address)
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


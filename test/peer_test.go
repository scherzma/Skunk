package test

import (
	"context"
	"testing"
	"time"

	"nhooyr.io/websocket"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/peer"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/tor"
	"github.com/stretchr/testify/assert"
)

const (
	waitTime = 1 * time.Second
)

func TestNewPeer(t *testing.T) {
	peerInstance, err := peer.NewPeer("127.0.0.1", "8080", "", "")
	assert.NoError(t, err)
	assert.NotNil(t, peerInstance)
	assert.Equal(t, peerInstance.Hostname, "127.0.0.1")
	assert.Equal(t, peerInstance.Port, "8080")
	assert.Equal(t, peerInstance.ProxyAddr, "")
}

func TestListen(t *testing.T) {
	peerInstance, err := peer.NewPeer("127.0.0.1", "8080", "", "")
	defer peerInstance.Shutdown()
	assert.NoError(t, err)

	peerInstance.Listen(nil)

	// Connecting to peer address should work
	conn, _, err := websocket.Dial(context.Background(), "ws://127.0.0.1:8080", nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	conn.Close(websocket.StatusNormalClosure, "test completed")
}

func TestConnect(t *testing.T) {
	peer1, err := peer.NewPeer("127.0.0.1", "1234", "", "")
	defer peer1.Shutdown()
	peer2, err := peer.NewPeer("127.0.0.1", "4321", "", "")
	defer peer2.Shutdown()

	peer1.Listen(nil)
	peer2.Listen(nil)

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
}

func TestPeerSetWriteConn(t *testing.T) {
	peer1, _ := peer.NewPeer("127.0.0.1", "1111", "", "")
	defer peer1.Shutdown()
	peer2, _ := peer.NewPeer("127.0.0.1", "10000", "", "")
	defer peer2.Shutdown()

	peer1.Listen(nil)

	err := peer2.Connect(peer1.Address)

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
}

func TestPeerReadMessages(t *testing.T) {
	peer1, _ := peer.NewPeer("127.0.0.1", "2222", "", "")
	defer peer1.Shutdown()
	peer2, _ := peer.NewPeer("127.0.0.1", "3333", "", "")
	defer peer2.Shutdown()
	peer3, _ := peer.NewPeer("127.0.0.1", "4444", "", "")
	defer peer3.Shutdown()
	peer4, _ := peer.NewPeer("127.0.0.1", "5555", "", "")
	defer peer4.Shutdown()
	peer5, _ := peer.NewPeer("127.0.0.1", "6666", "", "")
	defer peer5.Shutdown()

	peer1.Listen(nil)

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

	// Reading a message from each connection should work
	peer2.WriteMessage("Hello World!")
	peer3.WriteMessage("This is the story of my life")
	peer4.WriteMessage("Just do it!")
	peer5.WriteMessage("ABCDEFGHIJKLM")

	// Reading messages from the same connections again should work
	peer2.WriteMessage("Hello Proxima Centauri!")
	peer3.WriteMessage("Are you alright?")
	peer4.WriteMessage("No I'm all left!")
	peer5.WriteMessage("XYZ")

	// Reading multiple messages from the same connection should work
	peer2.WriteMessage("Recursive ...")
	peer2.WriteMessage("Recursive ...")
	peer2.WriteMessage("Recursive ...")
	peer2.WriteMessage("Recursive ...")

	// wait until all messages have been sent and received
	time.Sleep(13 * time.Second)

	// check that no error occured in this time
	select {
	case err := <-errorCh:
		assert.NoError(t, err)
	default:
	}
}

func TestPeerWriteMessage(t *testing.T) {
	peer1, _ := peer.NewPeer("127.0.0.1", "8888", "", "")
	defer peer1.Shutdown()
	peer2, _ := peer.NewPeer("127.0.0.1", "7890", "", "")
	defer peer2.Shutdown()

	peer1.Listen(nil)

	peer2.Connect(peer1.Address)
	peer2.SetWriteConn(peer1.Address)

	messageCh := make(chan string)
	errorCh := make(chan error)
	go peer1.ReadMessages(messageCh, errorCh)

	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{"numbers", "1234567890", "1234567890"},
		{"LETTERS", "ABCDEFGHIZ", "ABCDEFGHIZ"},
		{"letters", "abcdefghiz", "abcdefghiz"},
		{"special", "!?({&=$-:,", "!?({&=$-:,"},
		{"weird", "\t\n\r¬² ", "\t\n\r¬² "},
		{"mixture", "abc123ABC!", "abc123ABC!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			peer2.WriteMessage(tt.input)

			select {
			case msg := <-messageCh:
				assert.Equal(t, msg, tt.want)
			case err := <-errorCh:
				assert.NoError(t, err)
			}
		})
	}
}

func TestPeerShutdown(t *testing.T) {
	peerInstance, _ := peer.NewPeer("127.0.0.1", "1111", "", "")
	defer peerInstance.Shutdown()

	peerInstance.Listen(nil)

	peerInstance.Shutdown()

	// After shutdown you shoudn't be able to connect to the peer
	_, _, err := websocket.Dial(context.Background(), "ws://127.0.0.1:1111", nil)
	assert.Error(t, err)

	peerInstance.Listen(nil)

	// After executing Listen you should be able to connect to the peer again
	conn, _, err := websocket.Dial(context.Background(), "ws://127.0.0.1:1111", nil)
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	conn.Close(websocket.StatusNormalClosure, "test completed")
}

// this test represents two peers exchanging a message over the tor network
// because go-libtor can only start one tor process, we have to use one local installation of tor too
func TestPeerTor(t *testing.T) {
	// starts the non-embedded version of tor
	conf := &tor.TorConfig{
		DataDir:              "data-dir1",
		SocksPort:            "9052",
		LocalPort:            "1110",
		RemotePort:           "1111",
		DeleteDataDirOnClose: true,
		UseEmbedded:          false,
	}
	myTorOne, err := tor.NewTor(conf)
	assert.NoError(t, err)

	err = myTorOne.StartTor()
	assert.NoError(t, err)

	onionOne, err := myTorOne.StartHiddenService()
	assert.NoError(t, err)

	peerInstanceOne, _ := peer.NewPeer(onionOne.ID+".onion", "1110", "1111", "127.0.0.1:9052")
	defer peerInstanceOne.Shutdown()

	peerInstanceOne.Listen(onionOne)

	// starts the embedded version of tor
	conf = &tor.TorConfig{
		DataDir:              "data-dir2",
		SocksPort:            "9053",
		LocalPort:            "2221",
		RemotePort:           "2222",
		DeleteDataDirOnClose: true,
		UseEmbedded:          true,
	}
	myTorTwo, err := tor.NewTor(conf)
	assert.NoError(t, err)

	err = myTorTwo.StartTor()
	assert.NoError(t, err)

	onionTwo, err := myTorTwo.StartHiddenService()
	assert.NoError(t, err)

	peerInstanceTwo, _ := peer.NewPeer(onionTwo.ID+".onion", "2221", "2222", "127.0.0.1:9053")
	defer peerInstanceTwo.Shutdown()

	err = peerInstanceTwo.Connect(peerInstanceOne.Address)
	assert.NoError(t, err)

	peerInstanceTwo.SetWriteConn(peerInstanceOne.Address)

	messageCh := make(chan string)
	errorCh := make(chan error)

	go peerInstanceOne.ReadMessages(messageCh, errorCh)

	peerInstanceTwo.WriteMessage("Hello Dark World!")

	// wait until message has been sent and received
	time.Sleep(10 * time.Second)

	// check that no error occured during this time
	select {
	case msg := <-messageCh:
		t.Log(msg)
	case err := <-errorCh:
		assert.NoError(t, err)
	default:
	}

	myTorOne.StopTor()
	myTorTwo.StopTor()
}

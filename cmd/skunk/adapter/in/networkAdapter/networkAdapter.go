package networkAdapter

import (
	"encoding/json"
	"fmt"
	"sync"

	cretztor "github.com/cretz/bine/tor"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/peer"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/tor"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/util"
)

// Constants for configuration of Tor
const (
	SocksPort            = "9055"
	LocalPort            = "1111"
	RemotePort           = "2222"
	ReusePrivateKey      = true // reuse private key for constant onion address
	DeleteDataDirOnClose = false
	UseEmbedded          = true // use embedded tor process
)

var (
	networkAdapter *NetworkAdapter
	once           sync.Once
)

// NetworkAdapter connects the main logic to the tor peer network
type NetworkAdapter struct {
	subscriber network.NetworkObserver // subscriber observing network messages
	peer       *peer.Peer
	tor        *tor.Tor
}

func NewAdapter() *NetworkAdapter {
	// singleton
	once.Do(func() {
		networkAdapter = &NetworkAdapter{}
	})

	return networkAdapter
}

func (n *NetworkAdapter) SubscribeToNetwork(observer network.NetworkObserver) error {
	if n.subscriber == observer {
		return fmt.Errorf("network adapter is already connected to observer: %v", observer)
	}

	if n.peer != nil || n.tor != nil {
		return fmt.Errorf("network services are already running")
	}

	torInstance, onionService, err := startTor()
	if err != nil {
		return err
	}

	peerInstance, err := startPeer(onionService)
	if err != nil {
		return err
	}

	// begin asynchronously reading network messages.
	go n.readNetworkMessages()

	n.subscriber = observer
	n.peer = peerInstance
	n.tor = torInstance

	// notify the subscriber that the network is now online and send the onion address
	message := network.Message{
		Id:              util.UUID(),
		Timestamp:       util.CurrentTimeMillis(),
		Content:         fmt.Sprintf(`"%s"`, n.peer.Address),
		FromUser:        "",
		SenderAddress:   "",
		ReceiverAddress: "",
		ChatID:          "",
		Operation:       network.NETWORK_ONLINE,
	}
	n.SendNetworkMessageToSubscriber(message)
	return nil
}

// UnsubscribeFromNetwork stops all network services and unsubscribes the observer
func (n *NetworkAdapter) UnsubscribeFromNetwork() error {
	if n.subscriber == nil {
		return fmt.Errorf("can't unsubscribe from nil")
	}

	if n.peer == nil {
		return fmt.Errorf("peer network is nil")
	}

	if n.tor == nil {
		return fmt.Errorf("tor network is nil")
	}

	// stop tor and peer services
	n.tor.StopTor()
	n.peer.Shutdown()

	n.subscriber = nil
	n.peer = nil
	n.tor = nil
	return nil
}

// SendMessageToNetworkPeer sends a message to a specified network peer
func (n *NetworkAdapter) SendMessageToNetworkPeer(address string, message network.Message) error {
	// connect to peer if not already connected
	if !n.peer.IsConnectedTo(address) {
		err := n.peer.Connect(address)
		// if there is any error, we treat it as if the peer is offline
		if err != nil {
			// message subscriber that the peer is offline
			message := network.Message{
				Id:              util.UUID(),
				Timestamp:       util.CurrentTimeMillis(),
				Content:         "",
				FromUser:        "",
				SenderAddress:   n.peer.Address,
				ReceiverAddress: address,
				ChatID:          "",
				Operation:       network.USER_OFFLINE,
			}
			n.SendNetworkMessageToSubscriber(message)
			return nil
		}
	}

	err := n.peer.SetWriteConn(address)
	if err != nil {
		return err
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = n.peer.WriteMessage(string(jsonMessage))
	if err != nil {
		return err
	}
	return nil
}

// SendNetworkMessageToSubscriber forwards a network message to the subscribed observer
func (n *NetworkAdapter) SendNetworkMessageToSubscriber(message network.Message) {
	n.subscriber.Notify(message)
}

// readNetworkMessages reads messages from the network and forwards them to the subscriber
func (n *NetworkAdapter) readNetworkMessages() {
	messageCh := make(chan string)
	errorCh := make(chan error)

	// read messages asynchronously; terminate on unsubscribe
	// also closes messageCh and errorCh
	go n.peer.ReadMessages(messageCh, errorCh)

	// handle incoming messages and errors
	for {
		select {
		case msg, ok := <-messageCh:
			if !ok {
				return
			}
			message := network.Message{}
			err := json.Unmarshal([]byte(msg), &message)
			if err != nil {
				continue // optionally we could handle incorrect formatted messages
			}
			n.SendNetworkMessageToSubscriber(message)
		case err, ok := <-errorCh:
			if !ok {
				return
			}
			// message subscriber that a peer is offline
			message := network.Message{
				Id:              util.UUID(),
				Timestamp:       util.CurrentTimeMillis(),
				Content:         "",
				FromUser:        "",
				SenderAddress:   n.peer.Address,
				ReceiverAddress: err.Error(), // error contains the address
				ChatID:          "",
				Operation:       network.USER_OFFLINE,
			}
			n.SendNetworkMessageToSubscriber(message)
		}
	}
}

// startTor initializes and starts a Tor service.
func startTor() (*tor.Tor, *cretztor.OnionService, error) {
	conf := tor.TorConfig{
		SocksPort:            SocksPort,
		LocalPort:            LocalPort,
		RemotePort:           RemotePort,
		ReusePrivateKey:      ReusePrivateKey,
		DeleteDataDirOnClose: DeleteDataDirOnClose,
		UseEmbedded:          UseEmbedded,
	}

	torInstance, err := tor.NewTor(&conf)
	if err != nil {
		return nil, nil, err
	}

	err = torInstance.StartTor()
	if err != nil {
		return nil, nil, err
	}

	onionService, err := torInstance.StartHiddenService()
	if err != nil {
		return nil, nil, err
	}

	return torInstance, onionService, nil
}

// startPeer initializes a peer network connection using the provided OnionService
func startPeer(onionService *cretztor.OnionService) (*peer.Peer, error) {
	peerInstance, err := peer.NewPeer(onionService.ID+".onion", LocalPort, RemotePort, "127.0.0.1:"+SocksPort)
	if err != nil {
		return nil, err
	}
	peerInstance.Listen(onionService)

	return peerInstance, err
}

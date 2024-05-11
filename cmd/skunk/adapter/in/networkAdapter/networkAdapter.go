package networkAdapter

import (
	"encoding/json"
	"fmt"

	"github.com/cretz/bine/tor"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/peer"
	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/tor"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/util/timestamp"
	"github.com/scherzma/Skunk/cmd/skunk/util/uuid"
)

const (
	SocksPort            = "9055"
	LocalPort            = "1111"
	RemotePort           = "2222"
	DeleteDataDirOnClose = false
	UseEmbedded          = true
)

var (
	networkAdapter *NetworkAdapter
	once           sync.Once
)

type NetworkAdapter struct {
	subscriber *network.NetworkObserver
	peer       *peer.Peer
	tor        *tor.Tor
}

func NewAdapter() *NetworkAdapter {
	once.Do(func() {
		networkAdapter := NetworkAdapter{}
	})

	return &networkAdapter
}

func (n *NetworkAdapter) SubscribeToNetwork(observer *network.NetworkObserver) error {
	if n.subscriber == observer {
		return fmt.Errorf("network adapter is already connected to observer: %v", observer)
	}

    if n.peer != nil {
        return fmt.Errof("peer network is already running")
    }

    if n.tor != nil {
        return fmt.Errof("tor network is already running")
    }

    torInstance, onionService, err := startTor()
    if err != nil {
        return err
    }

    peerInstance, err := startPeer(onionService)
    if err != nil {
        return err
    }

    go readNetworkMessages()

    n.subscriber = observer
    n.peer = peerInstance
    n.tor = torInstance
	return nil
}

func (n *NetworkAdapter) UnsubscribeFromNetwork() error {
	if n.subscriber == nil {
		return fmt.Errorf("can't unsubscribe from nil")
	}

    if n.peer == nil {
        return fmt.Errof("peer network is nil")
    }

    if n.tor == nil {
        return fmt.Errof("tor network is nil")
    }

	stopTor()
	stopPeer()

	n.subscriber = nil
	n.peer = nil
	n.tor = nil
	return nil
}

func (n *NetworkAdapter) SendMessageToNetworkPeer(address string, message *network.Message) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	if !n.peer.IsConnectedTo(address) {
		err = n.peer.Connect(address)
        if err != nil {
            // message subscriber that the peer is offline
            message := network.Message {
                Id: uuid.UUID(),
                Timestamp: timestamp.CurrentTimeMillis(),
                Content: "",
                FromUser: "",
                SenderAddress: n.peer.Address,
                ReceiverAddress: address,
                ChatID: "",
                Operation: network.USER_OFFLINE
            }
            n.SendNetworkMessageToSubscriber(message)
            return nil
        }
	}

	err = n.peer.SetWriteConn(address) if err != nil {
		return err
	}

	err = n.peer.WriteMessage(string(jsonMessage))
	if err != nil {
		return err
	}
}

func (n *NetworkAdapter) SendNetworkMessageToSubscriber(message network.Message) {
	n.subscriber.Notify(message)
}

func startTor() (*tor.Tor, *tor.OnionService, error) {
	conf := tor.TorConfig{
		SocksPort:            SocksPort,
		LocalPort:            LocalPort,
		RemotePort:           RemotePort,
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

	return &torInstance, &onionService, nil
}

func startPeer(onionService *tor.OnionService) (*peer.Peer, error) {
	peerInstance, err := peer.NewPeer(onionService.ID+".onion", LocalPort, RemotePort, "127.0.0.1:"+SocksPort)
	if err != nil {
		return nil, err
	}

	return &peerInstance, err
}

func readNetworkMessages() {
	messageCh := make(chan string)
	errorCh := make(chan error)

	// this terminates on UnsubscribeFromNetwork
	// and closes messageCh and errorCh
	go n.peer.ReadMessages(messageCh, errorCh)

	for {
		select {
		case msg, ok := <-messageCh:
			if !ok {
				return
			}
			message := network.Message{}
			json.Unmarshal([]byte(msg), &message)

			n.SendNetworkMessageToSubscriber(message)
		case err, ok := <-errorCh:
			if !ok {
				return
			}
            // message subscriber that a peer is offline
            message := network.Message {
                Id: uuid.UUID(),
                Timestamp: timestamp.CurrentTimeMillis(),
                Content: "",
                FromUser: "",
                SenderAddress: n.peer.Address,
                ReceiverAddress: err.Error(),   // address is encoded in err
                ChatID: "",
                Operation: network.USER_OFFLINE
            }
			n.SendNetworkMessageToSubscriber(message)
		}
	}
}

func stopTor(torInstance *tor.Tor) {
	torInstance.StopTor()
}

func stopPeer(peerInstance *tor.Tor) {
	peerInstance.Shutdown()
}

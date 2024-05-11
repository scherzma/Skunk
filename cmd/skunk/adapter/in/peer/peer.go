package peer

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
    "github.com/pkg/errors"
	"golang.org/x/net/proxy"
)

const (
	MaxConns          = 64                           // MaxConns defines the maximum number of concurrent websocket connections allowed.
	connWait          = 1 * time.Minute              // connWait specifies the timeout for connecting to another peer.
	writeWait         = 20 * time.Second             // writeWait specifies the timeout for writing to another peer. has to be high when running over tor
	shutdownWait      = 0 * time.Second              // shutdownWait specifies the wait time for shutting down the HTTP server. (optional for later)
	readRateInterval  = 2 * time.Second              // readRateInterval specifies the rate at which it will it is tried to read a message from every connection.
	readWait          = (readRateInterval * 9) / 10  // readWait specifies the time for trying to read a message from a connection. Needs to be less than readRateInterval
	heartbeatInterval = 1 * time.Minute              // heartbeatInterval specifies the time interval between consecutive heartbeat messages. when running over tor this should be high
	pongWait          = (heartbeatInterval * 9) / 10 // pongWait specifies the time a ping response can need before the connection gets classified as closed during an heartbeat. Needs to be less than heartbeatInteval
	maxMessageSize    = 512                          // maxMessageSize defines the maximum message size allowed from peer. (bytes)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

// Peer encapsulates the state and functionality for a network peer, including its connections,
// configuration parameters, and synchronisation primitives for safe concurrent access.
type Peer struct {
	client     *http.Client               // client is used to make HTTP requests with a custom transport, supporting proxy configuration.
	readConns  map[string]*websocket.Conn // readConns maintains a map of active websocket connections for reading, indexed by the remote address. Note: Maybe we can later use a sync.Map
	mapRWLock  sync.RWMutex               // mapRWLock provides concurrent access control for readConns map.
	writeConn  *websocket.Conn            // writeConn is a dedicated websocket connection reserved for writing messages.
	readMutex  sync.Mutex                 // readMutex provides concurrent access control for ReadMessage.
	writeMutex sync.Mutex                 // writeMutex provides concurrent access control for WriteMessage.
	quitch     chan struct{}              // quitch is used to signal the shutdown process for the peer.
	Hostname   string                     // Hostname specifies the network address of the peer.
	Port       string                     // Port on which the peer listens for incoming connections.
	Address    string                     // Address specifies the complete websocket address: ws://Hostname:Port
	ProxyAddr  string                     // ProxyAddr specifies the address of SOCKS5 proxy, if used for connections.
}

// NewPeer initializes a new Peer instance with the given network settings.
// It also configures the peer's HTTP client for optimal proxy support.
// hostname needs to include .de, .onion...
func NewPeer(hostname string, localPort string, remotePort string, proxyAddr string) (*Peer, error) {
	if remotePort == "" {
		remotePort = localPort
	}
	transport, err := createTransport(proxyAddr) // Attempts to create an HTTP transport, optionally configured with a SOCKS5 proxy.
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	p := Peer{
		readConns: make(map[string]*websocket.Conn),
		Hostname:  hostname,
		Port:      localPort,
		Address:   fmt.Sprintf("ws://%s:%s", hostname, remotePort),
		ProxyAddr: proxyAddr,
		quitch:    make(chan struct{}),
		client: &http.Client{
			Transport: transport,
		},
	}

	// starts the heartbeat mechanism
	p.startHeartbeat()

	return &p, nil
}

// createTransport configures and returns an HTTP transport mechanism.
// If a proxy address is provided, it configures the transport to use a SOCKS5 proxy.
func createTransport(proxyAddr string) (*http.Transport, error) {
	if proxyAddr != "" {
		dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, nil)
		if err != nil {
			return nil, err
		}
		return &http.Transport{
			Dial: dialer.Dial,
		}, nil
	}
	return &http.Transport{}, nil
}

// Listen sets up an HTTP server and starts listening on the configured port for incoming websocket connections.
// It also starts a goroutine for graceful shutdown handling upon receiving a signal on the quitch channel.
// Optionally it uses the listener provided by tor.
func (p *Peer) Listen(l net.Listener) {
	select {
	case _, ok := <-p.quitch:
		if !ok {
			p.quitch = make(chan struct{}) // if closed, then reopen channel
		}
	default:
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.handler) // Registers the main handler.

	srv := &http.Server{
		Handler: mux,
	}

	go func() {
		var err error
		if l != nil {
			// Use the provided listener
			err = srv.Serve(l)
		} else {
			// Listen on the specified port if no listener is provided
			srv.Addr = ":" + p.Port
			err = srv.ListenAndServe()
		}
		if err != http.ErrServerClosed {
			log.Printf("HTTP server listen failed: %v", err)
		}
	}()

	// Shuts down server when quitch gets closed
	go func() {
		select {
		case <-p.quitch:
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

			if err := srv.Shutdown(shutdownCtx); err != nil {
				log.Printf("HTTP server shutdown failed: %v", err)
			}
			defer cancel()
		}
	}()
}

// SetWriteConn designates a specific websocket connection, identified by its address, as the sole connection for sending messages.
// It verifies that the peer is currently connected to the specified address before setting the connection.
func (p *Peer) SetWriteConn(address string) error {
	if len(p.readConns) == 0 {
		return fmt.Errorf("peer is not connected to any address")
	}

	if !p.IsConnectedTo(address) {
		return fmt.Errorf("peer is not connected to address: %s", address)
	}

	p.mapRWLock.RLock()
	p.writeConn = p.readConns[address]
	p.mapRWLock.RUnlock()
	return nil
}

// Connect establishes a new websocket connection to the specified address and adds it to the pool of read connections.
// It performs connection setup with a timeout and utilizes the configured HTTP client for the connection attempt.
func (p *Peer) Connect(address string) error {
	if address == "" {
		return fmt.Errorf("empty address is not a valid address")
	}

	if address == p.Address {
		return fmt.Errorf("can't connect to own address: %s", p.Address)
	}

	if p.IsConnectedTo(address) {
		return fmt.Errorf("peer is already connected to address: %s", address)
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: connWait,
	}

	if p.client.Transport != nil {
		if transport, ok := p.client.Transport.(*http.Transport); ok {
			dialer.NetDial = transport.Dial
		}
	}

	// Set peer address as header so that the other peer knows at which port we are listening
	headers := http.Header{}
	headers.Add("X-Peer-Address", p.Address)

	c, _, err := dialer.Dial(address, headers)

	if err != nil {
		return fmt.Errorf("failed to dial websocket: %v", err)
	}

	return p.handleNewConnection(c, address)
}

// ReadMessage attempts to read a single message from the specified websocket connection.
// It locks the readMutex to ensure exclusive access to the connection during the read operation.
func (p *Peer) readMessage(conn *websocket.Conn, address string) (string, error) {
	if conn == nil {
		return "", fmt.Errorf("invalid connection: connection is nil")
	}

	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(readWait))
	p.readMutex.Lock()
	_, messageBytes, err := conn.ReadMessage()
	p.readMutex.Unlock()
	if err != nil {
        // check if error is because connection is closed
        // => gets handeld as peer offline
		if p.checkConnIsClosed(address, err) {
			return "", err
        // every other error gets ignored
		} else {
			return "", nil
		}
	}

	return string(messageBytes), nil
}

// ReadMessages starts readMessage for every conn in readConns in the readRate interval.
func (p *Peer) ReadMessages(messageCh chan<- string, errorCh chan<- error) {
	ticker := time.NewTicker(readRateInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				p.mapRWLock.RLock()
				for addr, conn := range p.readConns {
					go func() {
						msg, err := p.readMessage(conn, addr)
						if err != nil {
                            // encode address as error so that we know which peer is offline
                            errOffline := errors.New(addr)
							errorCh <- errOffline
						} else if msg != "" { // "" can happen when an error occurs when reading from
							// the connection but the error is not due to the
							// connection being closed.
							messageCh <- msg
						}
					}()
				}
				p.mapRWLock.RUnlock()
			case <-p.quitch:
				ticker.Stop()
                close(messageCh)
                close(errorCh)
				return
			}
		}
	}()
}

// WriteMessage sends a message using the designated write connection.
// It locks the writeMutex to ensure exclusive access to the connection during the write operation.
func (p *Peer) WriteMessage(message string) error {
	if p.writeConn == nil {
		return fmt.Errorf("no write connection is set")
	}

	// Append From p.Address: to message so that the other peers knows which peer sent him the message.
	fullMessage := fmt.Sprintf("From %s: %s", p.Address, message)

	p.writeConn.SetWriteDeadline(time.Now().Add(writeWait))
	p.writeMutex.Lock()
	err := p.writeConn.WriteMessage(websocket.TextMessage, []byte(fullMessage))
	p.writeMutex.Unlock()
	if err != nil {
		p.checkConnIsClosed(p.writeConn.RemoteAddr().String(), err)
		return err
	}

	return nil
}

// Shutdown initiates the shutdown process for the peer, closing all active websocket connections.
// and signaling the quitch channel to stop the HTTP server.
func (p *Peer) Shutdown() {
	p.mapRWLock.RLock()
	defer p.mapRWLock.RUnlock()

	for _, conn := range p.readConns {
		conn.Close()
	}
	p.readConns = make(map[string]*websocket.Conn) // Resets the connection pool.
	p.writeConn = nil

	close(p.quitch) // Signals the shutdown listener to initiate server shutdown.
}

// IsConnectedTo checks if there is an existing websocket connection to the specified address.
func (p *Peer) IsConnectedTo(address string) bool {
	p.mapRWLock.RLock()
	_, ok := p.readConns[address]
	p.mapRWLock.RUnlock()
	return ok
}

// handler is the HTTP request handler for upgrading incoming requests to websocket connections.
// It accepts a websocket connection and adds it to the pool of read connections.
func (p *Peer) handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade incoming connection: %v", err)
		return
	}

	if err := p.handleNewConnection(conn, r.Header.Get("X-Peer-Address")); err != nil {
		log.Printf("failed to handle new connection: %v", err)
		conn.Close() // TO-DO: Send and handle error that the peer reached its maximum of connections.
	}
}

// handleNewConnection adds a newly established websocket connection to the readConns map.
// It ensures that the total number of connections does not exceed the maximum allowed.
func (p *Peer) handleNewConnection(conn *websocket.Conn, address string) error {
	if len(p.readConns) >= MaxConns {
		return fmt.Errorf("maximum number of connections reached: %d", MaxConns)
	}

	p.mapRWLock.Lock()
	p.readConns[address] = conn
	p.mapRWLock.Unlock()
	return nil
}

// StartHearbeat initiates a periodic hearbeat mechanism for all active connections.
// It sends a ping message at regular intervals to each connection to ensure they are alive.
func (p *Peer) startHeartbeat() {
	ticker := time.NewTicker(heartbeatInterval)
	go func() {
		for {
			select {
			case <-ticker.C: // On each tick, send heartbeat to all connections.
				p.sendHeartbeatToAll()
			case <-p.quitch:
				ticker.Stop()
				return
			}
		}
	}()
}

// sendHearbeatToAll sends a hearbeat signal (ping) to each active connection.
// If a connection fails to respond to the heartbeat, it removes the connection.
func (p *Peer) sendHeartbeatToAll() {
	//p.writeMutex.Lock()
	//defer p.writeMutex.Unlock()

	// THIS IS NOT GOOD
	for address, conn := range p.readConns {
		if conn == nil {
			p.mapRWLock.Lock()
			delete(p.readConns, address)
			p.mapRWLock.Unlock()
		} else {
			conn.SetWriteDeadline(time.Now().Add(pongWait))
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				p.mapRWLock.Lock()
				delete(p.readConns, address) // Optionally, we could try to reinitialize the connection here.
				p.mapRWLock.Unlock()
			}
		}
	}
}

// checkConnIsClosed evaluates if an error during a read or write operation was due to the connection being closed.
// If so, it removes the connection from the readConns map to prevent furter use.
func (p *Peer) checkConnIsClosed(address string, err error) bool {
	if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
		p.mapRWLock.Lock()
		delete(p.readConns, address)
		p.mapRWLock.Unlock()
		return true
	}
	return false
}

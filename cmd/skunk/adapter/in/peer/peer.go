package peer

import (
    "bytes"
	"context"
    "fmt"
	"log"
	"net/http"
    "sync"
	"time"

	"golang.org/x/net/proxy"
    "github.com/gorilla/websocket"
)

// TO-DO: Limit Message Rate (per Minute and message speed)

const (
    MaxConns = 64                           // MaxConns defines the maximum number of concurrent websocket connections allowed.
    connWait = 1 * time.Minute              // connWait specifies the timeout for connecting to another peer.
    writeWait = 10 * time.Second            // writeWait specifies the timeout for writing a heartbeat message.
    shutdownWait = 5 * time.Second              // shutdownWait specifies the wait time for shutting down the HTTP server.
    readRateInterval = 2 * time.Second
    heartbeatInterval = 2 * time.Minute    // heartbeatInterval specifies the time interval between consecutive hearbeat messages.
    maxMessageSize = 512                    // maxMessageSize defines the maximum message size allowed from peer.
)

var (
    newline = []byte{'\n'}
    space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

// Peer encapsulates the state and functionality for a network peer, including its connections,
// configuration parameters, and synchronisation primitives for safe concurrent access.
type Peer struct {
	client     *http.Client                 // client is used to make HTTP requests with a custom transport, supporting proxy configuration.
    readConns  map[string]*websocket.Conn   // readConns maintains a map of active websocket connections for reading, indexed by the remote address. Note: Maybe we can later use a sync.Map
    mapRWLock  sync.RWMutex                 // mapRWLock provides concurrent access control for readConns map.
	writeConn  *websocket.Conn              // writeConn is a dedicated websocket connection reserved for writing messages.
	Hostname   string                       // Hostname specifies the network address of the peer.
	Port       string                       // Port on which the peer listens for incoming connections.
    Address    string                       // Address specifies the complete websocket address: ws://Hostname:Port
	ProxyAddr  string                       // ProxyAddr specifies the address of SOCKS5 proxy, if used for connections.
    readMutex  sync.Mutex                   // readMutex provides concurrent access control for readConns.
    writeMutex sync.Mutex                   // writeMutex provides concurrent access control for writeConn.
	quitch     chan struct{}                // quitch is used to signal the shutdown process for the peer.
}

// NewPeer initializes a new Peer instance with the given network settings.
// It also configures the peer's HTTP client for optimal proxy support.
func NewPeer(hostname, port, proxyAddr string) (*Peer, error) {
	transport, err := createTransport(proxyAddr) // Attempts to create an HTTP transport, optionally configured with a SOCKS5 proxy.
	if err != nil {
        return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

    p := Peer{
        readConns: make(map[string]*websocket.Conn),
		Hostname:  hostname,
		Port:      port,
        Address:   fmt.Sprintf("ws://%s:%s", hostname, port),
		ProxyAddr: proxyAddr,
		quitch:    make(chan struct{}),
		client: &http.Client{
			Transport: transport,
		},
    }

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
func (p *Peer) Listen() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", p.handler) // Registers the main handler for incoming websocket upgrade requests.

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", p.Port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("HTTP server listen failed: %v", err)
		}
	}()

	go func() {
		for {
			_, closed := <-p.quitch
			if closed {
				shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownWait)
				defer cancel()

				if err := srv.Shutdown(shutdownCtx); err != nil {
					log.Printf("HTTP server shutdown failed: %v", err)
				}
				break
			}
		}
	}()
}

// SetWriteConn designates a specific websocket connection, identified by its address, as the sole connection for sending messages.
// It verifies that the peer is currently connected to the specified address before setting the connection.
func (p *Peer) SetWriteConn(address string) error {
    p.mapRWLock.RLock()
    defer p.mapRWLock.RUnlock()

    if len(p.readConns) == 0 {
         return fmt.Errorf("peer is not connected to any address")
    }

    if !p.isConnectedTo(address) {
        return fmt.Errorf("peer is not connected to address: %s", address) 
    }

    p.writeConn = p.readConns[address]
    return nil
}

// Connect establishes a new websocket connection to the specified address and adds it to the pool of read connections.
// It performs connection setup with a timeout and utilizes the configured HTTP client for the connection attempt.
func (p *Peer) Connect(address string) error {
    if p.isConnectedTo(address) {
        return fmt.Errorf("peer is already connected to address: %s", address)
    }

    dialer := websocket.Dialer{
        Proxy: http.ProxyFromEnvironment,
        HandshakeTimeout: connWait,
    }

    if p.client.Transport != nil {
        if transport, ok := p.client.Transport.(*http.Transport); ok {
            dialer.NetDial = transport.Dial
        }
    }

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
func (p *Peer) readMessage(conn *websocket.Conn) (string, error) {
    p.readMutex.Lock()
    defer p.readMutex.Unlock()

	if conn == nil {
		return "", fmt.Errorf("invalid connection: connection is nil")
	}

    conn.SetReadLimit(maxMessageSize)
    conn.SetReadDeadline(time.Now().Add(connWait))
    _, messageBytes, err := conn.ReadMessage()
    if err != nil {
        p.checkConnIsClosed(conn, err)
        return "", err
    }
    messageBytes = bytes.TrimSpace(bytes.Replace(messageBytes, newline, space, -1))

	return string(messageBytes), nil
}

func (p *Peer) ReadMessages(messageCh chan<- string, errorCh chan<- error) {
    ticker := time.NewTicker(readRateInterval)
    go func(){
        for {
            select{
                case <-ticker.C:
                    for _, conn := range p.readConns {
                        go func() {
                            msg, err := p.readMessage(conn)
                            if err != nil {
                                errorCh<-err
                            } else {
                                messageCh<-msg
                            }
                        }()
                    }
                case <-p.quitch:
                    ticker.Stop()
                    return
            }
        }
    }()
}

// WriteMessage sends a message using the designated write connection.
// It locks the writeMutex to ensure exclusive access to the connection during the write operation.
func (p *Peer) WriteMessage(message string) error {
    p.writeMutex.Lock()
    defer p.writeMutex.Unlock()

	if p.writeConn == nil {
		return fmt.Errorf("no write connection is set")
	}

    fullMessage := fmt.Sprintf("From %s: %s", p.Address, message)

    p.writeConn.SetWriteDeadline(time.Now().Add(writeWait))
    err := p.writeConn.WriteMessage(websocket.TextMessage, []byte(fullMessage))
    if err != nil {
        p.checkConnIsClosed(p.writeConn, err)
        return err
    }

    return nil
}

// Shutdown initiates the shutdown process for the peer, closing all active websocket connections.
// and signaling the quitch channel to stop the HTTP server.
func (p *Peer) Shutdown() {
    p.mapRWLock.RLock()
    for _, conn := range p.readConns {
        conn.Close()
    }
    p.readConns = make(map[string]*websocket.Conn) // Resets the connection pool.
    p.mapRWLock.RUnlock()

	close(p.quitch) // Signals the shutdown listener to initiate server shutdown.
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
    p.mapRWLock.Lock()
    defer p.mapRWLock.Unlock()

    if len(p.readConns) >= MaxConns {
        return fmt.Errorf("maximum number of connections reached: %d", MaxConns)
    }

    p.readConns[address] = conn
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
    p.writeMutex.Lock()
    defer p.writeMutex.Unlock()

    for address, conn := range p.readConns {
        if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
            p.mapRWLock.Lock()
            delete(p.readConns, address) // Optionally, we could try to reinitialize the connection here.
            p.mapRWLock.Unlock()
        }
    }
}

// isConnectedTo checks if there is an existing websocket connection to the specified address.
func (p *Peer) isConnectedTo(address string) bool {
    p.mapRWLock.RLock()
    defer p.mapRWLock.RUnlock()
    _, ok := p.readConns[address]
    return ok
}

// checkConnIsClosed evaluates if an error during a read or write operation was due to the connection being closed.
// If so, it removes the connection from the readConns map to prevent furter use.
func (p *Peer) checkConnIsClosed(conn *websocket.Conn, err error) {
	if conn != nil && websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
        p.mapRWLock.Lock()
        delete(p.readConns, conn.RemoteAddr().String())
        p.mapRWLock.Unlock()
	} else if conn == nil {
        p.mapRWLock.Lock()
        for addr, c := range p.readConns {
            if conn == c {
                delete(p.readConns, addr)
                break
            }
        }
        p.mapRWLock.Unlock()
    }
}


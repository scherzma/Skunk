package peer

import (
	"context"
	"fmt"
	"log"
	"net/http"
    "sync"
	"time"

	"golang.org/x/net/proxy"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const (
    MaxConns = 64                           // MaxConns defines the maximum number of concurrent websocket connections allowed.
    connWait = 1 * time.Minute              // connWait specifies the timeout for connecting to another peer.
    shutdownWait = 5 * time.Second              // shutdownWait specifies the wait time for shutting down the HTTP server.
    heartbeatInterval = 30 * time.Second    // heartbeatInterval specifies the time interval between consecutive hearbeat messages.
    writeWait = 10 * time.Second            // writeWait specifies the timeout for writing a heartbeat message.
)

// Peer encapsulates the state and functionality for a network peer, including its connections,
// configuration parameters, and synchronisation primitives for safe concurrent access.
type Peer struct {
	client     *http.Client                 // client is used to make HTTP requests with a custom transport, supporting proxy configuration.
    readConns  map[string]*websocket.Conn   // readConns maintains a map of active websocket connections for reading, indexed by the remote address.
    mapRWLock  sync.RWMutex                 // mapRWLock provides concurrent access control for readConns map.
	writeConn  *websocket.Conn              // writeConn is a dedicated websocket connection reserved for writing messages.
	Hostname   string                       // Hostname specifies the network address of the peer.
	Port       string                       // Port on which the peer listens for incoming connections.
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

	ctx, cancel := context.WithTimeout(context.Background(), connWait)
	defer cancel()

	c, _, err := websocket.Dial(ctx, address, &websocket.DialOptions{HTTPClient: p.client})
	if err != nil {
		return fmt.Errorf("failed to dial websocket: %v", err)
	}

	return p.handleNewConnection(c)
}

// ReadMessage attempts to read a single message from the specified websocket connection.
// It locks the readMutex to ensure exclusive access to the connection during the read operation.
func (p *Peer) ReadMessage(address string) (interface{}, error) {
    p.readMutex.Lock()
    defer p.readMutex.Unlock()

    if !p.isConnectedTo(address) {
        return nil, fmt.Errorf("not conected to %s", address)
    }

    p.mapRWLock.RLock()
    conn := p.readConns[address]
    p.mapRWLock.RUnlock()

	if conn == nil {
		return nil, fmt.Errorf("invalid connection: connection is nil")
	}

	var msg interface{}
	err := wsjson.Read(context.Background(), conn, &msg)
    if err != nil {
        p.checkConnIsClosed(conn, err) // Verifies if the connection should be removed from the pool due to being closed.
        return nil, err
    }

	return msg, nil
}

// WriteMessage sends a message using the designated write connection.
// It locks the writeMutex to ensure exclusive access to the connection during the write operation.
func (p *Peer) WriteMessage(message string) error {
    p.writeMutex.Lock()
    defer p.writeMutex.Unlock()

	if p.writeConn == nil {
		return fmt.Errorf("no write connection is set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), writeWait)
	defer cancel()

	err := wsjson.Write(ctx, p.writeConn, message)
    if err != nil {
        p.checkConnIsClosed(p.writeConn, err) // Checks if the write connection is closed and updates state accordingly.
        return err
    }

    return nil
}

// Shutdown initiates the shutdown process for the peer, closing all active websocket connections.
// and signaling the quitch channel to stop the HTTP server.
func (p *Peer) Shutdown() {
    p.mapRWLock.RLock()
    for _, conn := range p.readConns {
        conn.Close(websocket.StatusNormalClosure, "")
    }
    p.readConns = make(map[string]*websocket.Conn) // Resets the connection pool.
    p.mapRWLock.RUnlock()

	close(p.quitch) // Signals the shutdown listener to initiate server shutdown.
}

// handler is the HTTP request handler for upgrading incoming requests to websocket connections.
// It accepts a websocket connection and adds it to the pool of read connections.
func (p *Peer) handler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Printf("failed to upgrade incoming connection: %v", err)
		return
	}

    if err := p.handleNewConnection(c); err != nil {
        log.Printf("failed to handle new connection: %v", err)
        c.Close(websocket.StatusNormalClosure, "") // TO-DO: Send and handle error that the peer reached its maximum of connections.
    }
}

// handleNewConnection adds a newly established websocket connection to the readConns map.
// It ensures that the total number of connections does not exceed the maximum allowed.
func (p *Peer) handleNewConnection(conn *websocket.Conn) error {
    p.mapRWLock.Lock()
    defer p.mapRWLock.Unlock()

    if len(p.readConns) >= MaxConns {
        return fmt.Errorf("maximum number of connections reached: %d", MaxConns)
    }
    p.readConns[conn.RemoteAddr().String()] = conn
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
        ctx, cancel := context.WithTimeout(context.Background(), writeWait)
        defer cancel()

        err := conn.Ping(ctx)
        if err != nil {
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
	if conn != nil && websocket.CloseStatus(err) != -1 {
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


package peer

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/proxy"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type Peer struct {
	client    *http.Client
	conn      *websocket.Conn
	Hostname  string
	Port      string
	ProxyAddr string
	quitch    chan struct{}
}

func NewPeer(hostname, port, proxyAddr string) (*Peer, error) {
	transport, err := createTransport(proxyAddr)
	if err != nil {
		return nil, fmt.Errorf("creating SOCKS5 dialer: %w", err)
	}

	return &Peer{
		Hostname:  hostname,
		Port:      port,
		ProxyAddr: proxyAddr,
		quitch:    make(chan struct{}),
		client: &http.Client{
			Transport: transport,
		},
	}, nil
}

func createTransport(proxyAddr string) (*http.Transport, error) {
	if proxyAddr != "" {
		dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, nil)
		if err != nil {
			return nil, fmt.Errorf("creating SOCKS5 dialer: %w", err)
		}
		return &http.Transport{
			Dial: dialer.Dial,
		}, nil
	} else {
		return &http.Transport{}, nil
	}
}

func (p *Peer) Listen() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", p.handler)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", p.Port),
		Handler: mux,
	}

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	go func() {
		for {
			_, closed := <-p.quitch
			if closed {
				shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), time.Second*5)
				defer shutdownRelease()

				if err := srv.Shutdown(shutdownCtx); err != nil {
					log.Fatal(err)
				}
				break
			}
		}
	}()
}

func (p *Peer) Connect(address string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, address, &websocket.DialOptions{HTTPClient: p.client})
	if err != nil {
		return err
	}

	p.closeCurrentConn()

	p.conn = c
	return nil
}

func (p *Peer) ReadMessage() (interface{}, error) {
	if p.conn == nil {
		return nil, fmt.Errorf("Peer needs to be connected to read from a connection")
	}
	var v interface{}
	err := wsjson.Read(context.Background(), p.conn, &v)

	p.checkConnIsClosed(err)
	return v, err
}

func (p *Peer) WriteMessage(message string) error {
	if p.conn == nil {
		return fmt.Errorf("Peer needs to be connected to write to a connection")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	err := wsjson.Write(ctx, p.conn, message)

	p.checkConnIsClosed(err)
	return err
}

func (p *Peer) Shutdown() {
	p.conn.Close(websocket.StatusNormalClosure, "")
	p.conn = nil
	close(p.quitch)
}

func (p *Peer) handler(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	p.closeCurrentConn()
	p.conn = c
}

func (p *Peer) closeCurrentConn() {
	if p.conn != nil {
		p.conn.Close(websocket.StatusNormalClosure, "")
		p.conn = nil
	}
}

func (p *Peer) checkConnIsClosed(err error) {
	if websocket.CloseStatus(err) != -1 {
		p.conn = nil
	}
}

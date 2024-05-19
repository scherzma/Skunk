package tor

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cretz/bine/tor"
	"github.com/ipsn/go-libtor"
)

// TorConfig holds configuration parameters for a Tor instance.
type TorConfig struct {
	DataDir              string // directory for storing tor data files.
	SocksPort            string // SOCKS5 proxy port for tor connection.
	LocalPort            string // local port for incoming connections.
	RemotePort           string // remote port for hiddenservice.
	DeleteDataDirOnClose bool   // flag to delete data directory on closing tor.
	ReusePrivateKey      bool   // flag to reuse the private key across sessions.
	UseEmbedded          bool   // use the embedded tor process (go-libtor)
}

// Tor wraps a tor configuration and instance for managing tor services.
type Tor struct {
	torConfig   *TorConfig // configuration for the tor instance.
	torInstance *tor.Tor   // the tor instance.
}

// NewTor initializes a new tor instance with the provided configuration.
func NewTor(torConfig *TorConfig) (*Tor, error) {
	// basic validation of configuration paramters.
	if torConfig.SocksPort == "" {
		return nil, fmt.Errorf("no tor socks port provided")
	}
	if torConfig.SocksPort == "9050" {
		return nil, fmt.Errorf("can't use 9050 as tor socks port")
	}
	if torConfig.LocalPort == "" {
		return nil, fmt.Errorf("no local port given")
	}
	if torConfig.RemotePort == "" {
		return nil, fmt.Errorf("no remote port given")
	}
	if torConfig.LocalPort == torConfig.SocksPort || torConfig.RemotePort == torConfig.SocksPort {
		return nil, fmt.Errorf("ports for hiddenservice can't match tor socks port")
	}

	if torConfig.DataDir == "" {
		torConfig.DataDir = "tor-data"
	}

	return &Tor{
		torConfig:   torConfig,
		torInstance: nil,
        onion: nil,
	}, nil
}

// StarTor starts the tor instance with the configured settings.
func (t *Tor) StartTor() error {
	if t.torInstance != nil {
		return fmt.Errorf("can't start the same tor instance twice")
	}

	conf := &tor.StartConf{
		NoAutoSocksPort: true, // needs to be true to be able to set a custom socks port
		ExtraArgs:       []string{"--SocksPort", t.torConfig.SocksPort},
		DataDir:         t.torConfig.DataDir,
		DebugWriter:     os.Stdout, // just for testing. might change later
	}

	// if configured use the go-libtor embedded tor process creator
	if t.torConfig.UseEmbedded {
		conf.ProcessCreator = libtor.Creator
	}

	torInstance, err := tor.Start(nil, conf)
	if err != nil {
		return err
	}

	t.torInstance = torInstance
	return nil
}

// StartHiddenService starts a hidden service using the current tor instance.
func (t *Tor) StartHiddenService() (*tor.OnionService, error) {
	if t.torInstance == nil {
		return nil, fmt.Errorf("tor needs to be started before a hiddenservice can be created")
	}

	// LocalPort and RemotePort are strings because we want to provide a unified interface for the tor config
	// where you don't need to figure out which port / address has which type.
	remotePortInt, err := strconv.Atoi(t.torConfig.RemotePort)
	if err != nil {
		return nil, err
	}
	localPortInt, err := strconv.Atoi(t.torConfig.LocalPort)
	if err != nil {
		return nil, err
	}

	conf := &tor.ListenConf{
		Version3:    true, // uses v3 onion service and ed25519 key
		LocalPort:   localPortInt,
		RemotePorts: []int{remotePortInt},
	}

	// attempt to read the private key if ReusePrivateKey is true
	if t.torConfig.ReusePrivateKey {
		privateKeyPath := filepath.Join(t.torConfig.DataDir, "hidden_service_private_key")
		var privateKey ed25519.PrivateKey

		if _, err := os.Stat(privateKeyPath); err == nil {
			keyData, readErr := os.ReadFile(privateKeyPath)
			if readErr != nil {
				return nil, fmt.Errorf("failed to read private key: %v", readErr)
			}
			privateKey = ed25519.PrivateKey(keyData)
		} else {
			// generate a new private key
			_, privKey, genErr := ed25519.GenerateKey(rand.Reader)
			if genErr != nil {
				return nil, fmt.Errorf("failed to generate private key: %v", genErr)
			}
			privateKey = privKey
			// save newly generated private key
			if err := os.WriteFile(privateKeyPath, privKey.Seed(), 0600); err != nil {
				return nil, fmt.Errorf("failed to save private key: %v", err)
			}
		}
		conf.Key = privateKey
	}

	// wait at most a few minutes to publish the service
	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()

	onion, err := t.torInstance.Listen(listenCtx, conf)
	if err != nil {
		return nil, err
	}
    t.onion = onion

	return onion, nil
}

// StopTor stops the tor instance (and hiddenservice) and handles cleanup.
func (t *Tor) StopTor() error {
    if t.onion == nil {
		return fmt.Errorf("need to wait until Listen returned before attempting to close")
    }
    if t.torInstance == nil {
		return fmt.Errorf("can't stop tor if tor is nil")
    }

	err := t.torInstance.Close()
	if err != nil {
		return fmt.Errorf("error stopping tor: %v", err)
	}

	// clean up data directory if configured
	if t.torConfig.DataDir != "" && t.torConfig.DeleteDataDirOnClose {
		if err := os.RemoveAll(t.torConfig.DataDir); err != nil {
			return fmt.Errorf("failed to remove data dir %v: %v", t.torConfig.DataDir, err)
		}
	}

	t.torInstance = nil
    t.onion = nil
	return nil
}

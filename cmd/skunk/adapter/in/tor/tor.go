package tor

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/cretz/bine/tor"
    "github.com/ipsn/go-libtor"
)

func StartTor(socksPort string, dataDir string, useEmbedded bool) (*tor.Tor, error) {
	if dataDir == "" {
		dataDir = "tor-data"
	}

	torConfig := &tor.StartConf{
		NoAutoSocksPort: true,
		// Should work without setting the SocksPolicy ("--SocksPolicy", "accept 127.0.0.1" 
        // To-Do: Test from different devices!!
		ExtraArgs: []string{"--SocksPort", socksPort},
		DataDir:   dataDir,
        DebugWriter: os.Stdout, // Note: Just for testing, change later
	}

    // You're only able to run one go-libtor instance, so for testing you need to be able to disable it
    if useEmbedded {
        torConfig.ProcessCreator = libtor.Creator
    }

	t, err := tor.Start(nil, torConfig)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func StartHiddenService(t *tor.Tor, localPort string, remotePort string) (string, *tor.OnionService, error) {
	// privateKey, err := getPrivateKey(t)
	// if err != nil {
	//     return "", nil, err
	// }

	remotePortInt, err := strconv.Atoi(remotePort)
	if err != nil {
		return "", nil, err
	}
	localPortInt, err := strconv.Atoi(localPort)

	listenCtx, listenCancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer listenCancel()

	// onion, err := t.Listen(listenCtx, &tor.ListenConf{Version3: true, Key: privateKey, LocalPort: localPortInt, RemotePorts: []int{remotePortInt}})
	onion, err := t.Listen(listenCtx, &tor.ListenConf{Version3: true, LocalPort: localPortInt, RemotePorts: []int{remotePortInt}})
	if err != nil {
		return "", nil, err
	}

	return onion.ID, onion, err
}

func StopTor(t *tor.Tor) {
	t.Close()
}

/*
func saveServiceInfo(key string) error {
    file, err := os.Create("serviceinfo.json")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    return encoder.Encode(key)
}

func loadServiceInfo() (string, error) {
    file, err := os.Open("serviceinfo.json")
    if err != nil {
        return "", err
    }
    defer file.Close()
    decoder := json.NewDecoder(file)
    var key string
    if err := decoder.Decode(&key); err != nil {
        return "", err
    }
    return key, nil
}

func getPrivateKey(t *tor.Tor) (ed25519.PrivateKey, error) {
    path := fmt.Sprintf("%s/serviceinfo.json", t.DataDir)

    if _, err := os.Stat(path); err == nil {
        privateKeyString, err := loadServiceInfo()
        if err != nil {
            return nil, err
        }

        privateKeyBytes, _ := hex.DecodeString(privateKeyString)
        privateKey := ed25519.PrivateKey(privateKeyBytes)

        return privateKey, nil
    } else {
        keyPair, _ := ed25519.GenerateKey(nil)
        privateKey := hex.EncodeToString(keyPair.PrivateKey())

        err = saveServiceInfo(privateKey)
        if err != nil {
            return nil, err
        }

        return keyPair.PrivateKey(), nil
    }
}
*/

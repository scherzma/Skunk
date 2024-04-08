package tor

import (
    "context"
    "encoding/hex"
    "encoding/json"
    "github.com/cretz/bine/torutil/ed25519"
    "os"
    "time"

    "github.com/cretz/bine/tor"
)

func StartTor() (*tor.Tor, error) {
    t, err := tor.Start(nil, &tor.StartConf{DataDir: "tor-data"})
    if err != nil {
        return nil, err
    }

    return t, nil
}

func StartHiddenService(t *tor.Tor) (string, error) {
    privateKey, err := getPrivateKey()

    listenCtx, listenCancel := context.WithTimeout(context.Background(), 1*time.Minute)
    defer listenCancel()

    onion, err := t.Listen(listenCtx, &tor.ListenConf{Version3: true, Key: privateKey, RemotePorts: []int{80}})
    if err != nil {
        return "", err
    }

    return onion.ID, err
}

func StopHiddenService(t *tor.Tor) {
    t.Close()
}

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

func getPrivateKey() (ed25519.PrivateKey, error) {
    if _, err := os.Stat("serviceinfo.json"); err == nil {
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

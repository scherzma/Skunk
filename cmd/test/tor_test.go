package test

import (
    "testing"

    "github.com/scherzma/Skunk/cmd/skunk/adapter/in/tor"
)

func TestTor(t *testing.T) {
    torInstance, _ := tor.StartTor()
    onionID, _ := tor.StartHiddenService(torInstance)

    t.Log(onionID)

    tor.StopHiddenService(torInstance)
}



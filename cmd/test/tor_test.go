package test

import (
    "testing"

    "github.com/scherzma/Skunk/cmd/skunk/adapter/in/tor"
)

func TestTor(t *testing.T) {
    torInstance, _ := tor.StartTor("9080", "")
    onionID, _, _ := tor.StartHiddenService(torInstance, "1111", "2222")

    t.Log(onionID)

    tor.StopTor(torInstance)
}


package test

import (
	"os"
	"testing"
	"time"

	"github.com/scherzma/Skunk/cmd/skunk/adapter/in/tor"
	"github.com/stretchr/testify/assert"
)

// The tests always use the embedded version of go
// since this will also be the version we'll use in production

func TestStartTor(t *testing.T) {
	conf := &tor.TorConfig{
		SocksPort:            "9070",
		LocalPort:            "1111",
		RemotePort:           "2222",
		DeleteDataDirOnClose: true,
		UseEmbedded:          true,
	}
	myTor, err := tor.NewTor(conf)
	assert.NoError(t, err)

	err = myTor.StartTor()
	assert.NoError(t, err)

	err = myTor.StopTor()
	assert.NoError(t, err)
}

func TestStartHiddenService(t *testing.T) {
	conf := &tor.TorConfig{
		SocksPort:            "9080",
		LocalPort:            "3333",
		RemotePort:           "4444",
		DeleteDataDirOnClose: true,
		UseEmbedded:          true,
	}
	myTor, err := tor.NewTor(conf)
	assert.NoError(t, err)

	err = myTor.StartTor()
	assert.NoError(t, err)

	onion, err := myTor.StartHiddenService()
	assert.NoError(t, err)
	if assert.NotNil(t, onion) {
		assert.NotNil(t, onion.LocalListener)
		assert.NotNil(t, onion.RemotePorts)
		assert.NotNil(t, onion.ID)
	}

	err = myTor.StopTor()
	assert.NoError(t, err)
}

func TestStopTor(t *testing.T) {
	dataDir := "test-data-dir"
	conf := &tor.TorConfig{
		DataDir:              dataDir,
		SocksPort:            "9090",
		LocalPort:            "5555",
		RemotePort:           "6666",
		DeleteDataDirOnClose: true,
		UseEmbedded:          true,
	}
	myTor, err := tor.NewTor(conf)
	assert.NoError(t, err)

	err = myTor.StartTor()
	assert.NoError(t, err)

	err = myTor.StopTor()
	assert.NoError(t, err)

	// the dataDir-folder shouldn't exist anymore
	_, err = os.Stat(dataDir)
	assert.True(t, os.IsNotExist(err))
}

func TestReusePrivateKeyTor(t *testing.T) {
	dataDir := "reuse-data-dir"
	// configure tor to reuse the private key
	conf := &tor.TorConfig{
		DataDir:              dataDir,
		SocksPort:            "9100",
		LocalPort:            "1234",
		RemotePort:           "4321",
		DeleteDataDirOnClose: false,
		ReusePrivateKey:      true,
		UseEmbedded:          true,
	}
	myTor, err := tor.NewTor(conf)
	assert.NoError(t, err)

	err = myTor.StartTor()
	assert.NoError(t, err)

	onionOne, err := myTor.StartHiddenService()
	assert.NoError(t, err)

	onionIDOne := onionOne.ID

	err = myTor.StopTor()
	assert.NoError(t, err)

	time.Sleep(5 * time.Second)

	// restart tor hidden service
	// and check if onionID is the same as before
	conf = &tor.TorConfig{
		DataDir:              dataDir,
		SocksPort:            "9100",
		LocalPort:            "5678",
		RemotePort:           "1010",
		DeleteDataDirOnClose: true,
		ReusePrivateKey:      true,
		UseEmbedded:          true,
	}
	myTor, err = tor.NewTor(conf)
	assert.NoError(t, err)

	err = myTor.StartTor()
	assert.NoError(t, err)

	onionTwo, err := myTor.StartHiddenService()
	assert.NoError(t, err)

	onionIDTwo := onionTwo.ID

	err = myTor.StopTor()
	assert.NoError(t, err)

	assert.Equal(t, onionIDOne, onionIDTwo)
}

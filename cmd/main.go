package main

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage/storageSQLiteAdapter"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
)

/*
type Message struct {
	Id              string
	Timestamp       int64
	Content         string
	SenderID        string
	ReceiverID      string
	SenderAddress   string
	ReceiverAddress string
	ChatID          string
	Operation       OperationType
}
*/

func main() {

	db := storageSQLiteAdapter.NewStorageSQLiteAdapter("test.db")
	db.CreateTables()
	db.StoreMessage(network.Message{
		Id:              "2",
		Timestamp:       1633029445,
		Content:         "Hello World!",
		SenderID:        "asda",
		ReceiverID:      "asdf4",
		SenderAddress:   "asdf",
		ReceiverAddress: "asdf",
		ChatID:          "asdf",
		Operation:       network.TEST_MESSAGE,
	})
	db.InvitedToChat("2", []storageSQLiteAdapter.PublicKeyAddress{
		{Address: "asdf", PublicKey: "asdf"},
		{Address: "asdf", PublicKey: "asdf"},
		{Address: "asdf", PublicKey: "asdf"},
	})

}

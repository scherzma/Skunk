package main

import (
	"github.com/scherzma/Skunk/cmd/skunk/adapter/out/storage"
)

func main() {
	// Create a new SQLite storage for the message queue
	storeMessageQueueSQLite := storage.StoreMessageQueueSQLite{}
	storeMessageQueueSQLite.StoreMessageQueue(nil)
}

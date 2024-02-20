package storage

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/scherzma/Skunk/cmd/skunk/application/port"
	"log"
)

type StoreMessageQueueSQLite struct {
}

func (s *StoreMessageQueueSQLite) StoreMessageQueue(messageQueue port.StoreMessageQueue) error {
	db, err := sql.Open("sqlite3", "./yourdatabase.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Check for table existence and create if not exists
	createTableSQL := `CREATE TABLE IF NOT EXISTS exampleTable (
        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,     
        "data" TEXT
    );` // SQL Statement for Create Table

	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *StoreMessageQueueSQLite) RetriveMessageQueue() (string, error) {
	return "", nil
}

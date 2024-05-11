package storageSQLiteAdapter

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
	"log"
)

type StorageSQLiteAdapter struct {
	db *sql.DB

	/* Check for table existence and create if not exists
		createTableSQL := `CREATE TABLE IF NOT EXISTS exampleTable (
	        "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	        "data" TEXT
	    );` // SQL Statement for Create Table

		_, err = db.Exec(createTableSQL)
		if err != nil {
		log.Fatal(err)
		}
	*/
}

func NewStorageSQLiteAdapter(db *sql.DB) *StorageSQLiteAdapter {
	db, err := sql.Open("sqlite3", "./skunk.db")
	if err != nil {
		log.Fatal(err)
	}
	return &StorageSQLiteAdapter{db: db}
}

func (a *StorageSQLiteAdapter) PeerSetUsername(peerId string, chatId string, username string) error {
	// Implement the logic to store the peer's username in the database
	// ...
	setUsernameSQL := `INSERT INTO peers (username, chat_id) VALUES (?, ?)`
	_, err := a.db.Exec(setUsernameSQL, username, chatId)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (a *StorageSQLiteAdapter) PeerJoinedChat(peerId string, chatId string) error {
	// Implement the logic to store the peer joining a chat in the database
	// ...
	return nil
}

func (a *StorageSQLiteAdapter) PeerLeftChat(peerId string, chatId string, senderId string, date int64) error {
	// Implement the logic to store the peer leaving a chat in the database
	// ...
	return nil
}

func (a *StorageSQLiteAdapter) ChatCreated(chatName string, chatId string) error {
	// Implement the logic to store a newly created chat in the database
	// ...
	return nil
}

func (a *StorageSQLiteAdapter) InvitatedToChat(chatId string, chatName string, peers []string) error {
	// Implement the logic to store a chat invitation in the database
	// ...
	return nil
}

func (a *StorageSQLiteAdapter) PeerGotInvitedToChat(peerId string, chatId string) error {
	// Implement the logic to store a peer getting invited to a chat in the database
	// ...
	return nil
}

func (a *StorageSQLiteAdapter) GetInvitations(peerId string) []string {
	// Implement the logic to retrieve a peer's chat invitations from the database
	// ...
	return []string{}
}

func (a *StorageSQLiteAdapter) GetMissingInternalMessages(chatId string, inputMessageIDs []string) []string {
	// Implement the logic to retrieve missing internal messages from the database
	// ...
	return []string{}
}

func (a *StorageSQLiteAdapter) GetMissingExternalMessages(chatId string, inputMessageIDs []string) []string {
	// Implement the logic to retrieve missing external messages from the database
	// ...
	return []string{}
}

func (a *StorageSQLiteAdapter) StoreMessage(message network.Message) error {
	// Implement the logic to store a network message in the database
	// ...
	return nil
}

func (a *StorageSQLiteAdapter) RetrieveMessage(messageId string) (network.Message, error) {
	// Implement the logic to retrieve a network message from the database
	// ...
	return network.Message{}, nil
}

func (a *StorageSQLiteAdapter) GetChats() []store.Chat {
	// Implement the logic to retrieve chats from the database
	// ...
	return []store.Chat{}
}

func (a *StorageSQLiteAdapter) GetUsername(peerId string) string {
	// Implement the logic to retrieve a peer's username from the database
	// ...
	return ""
}

func (a *StorageSQLiteAdapter) GetUsersInChat(chatId string) []store.User {
	// Implement the logic to retrieve users in a chat from the database
	// ...
	return []store.User{}
}

func (a *StorageSQLiteAdapter) GetPeers() []string {
	// Implement the logic to retrieve peers from the database
	// ...
	return []string{}
}

func (a *StorageSQLiteAdapter) GetChatMessages(chatId string) []store.ChatMessage {
	// Implement the logic to retrieve chat messages from the database
	// ...
	return []store.ChatMessage{}
}

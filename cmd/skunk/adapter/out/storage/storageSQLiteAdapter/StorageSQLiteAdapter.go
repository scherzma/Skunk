package storageSQLiteAdapter

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/store"
	"log"
	"strings"
	"sync"
)

type StorageSQLiteAdapter struct {
	db *sql.DB
}

var (
	instance *StorageSQLiteAdapter
	once     sync.Once
)

// NewStorageSQLiteAdapter creates a new instance of StorageSQLiteAdapter and initializes the database
func newStorageSQLiteAdapter(dbPath string) *StorageSQLiteAdapter {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	adapter := &StorageSQLiteAdapter{db: db}
	adapter.createTables()

	return adapter
}

// GetInstance returns the singleton instance of StorageSQLiteAdapter
func GetInstance(dbPath string) *StorageSQLiteAdapter {
	once.Do(func() {
		instance = newStorageSQLiteAdapter(dbPath)
	})
	return instance
}

func (a *StorageSQLiteAdapter) createTables() {
	sqlCommands := "CREATE TABLE IF NOT EXISTS Chats (\n    chat_id VARCHAR(1024) NOT NULL CONSTRAINT Chats_pk PRIMARY KEY,\n    name VARCHAR(40) NOT NULL\n);\n\nCREATE TABLE IF NOT EXISTS Peers (\n    peer_id INTEGER NOT NULL CONSTRAINT Peers_pk PRIMARY KEY AUTOINCREMENT,\n    public_key VARCHAR(1024) NOT NULL,\n    address VARCHAR(1024) NOT NULL\n);\n\nCREATE TABLE IF NOT EXISTS ChatMembers (\n    chat_member_id INTEGER NOT NULL CONSTRAINT ChatMembers_pk PRIMARY KEY AUTOINCREMENT,\n    date INTEGER NOT NULL,\n    peer_id VARCHAR(1024) NOT NULL CONSTRAINT ChatMembers_Peers_peer_id_fk REFERENCES Peers ON UPDATE CASCADE,\n    chat_id VARCHAR(1024) NOT NULL CONSTRAINT ChatMembers_Chats_chat_id_fk REFERENCES Chats ON UPDATE CASCADE,\n    username VARCHAR(50)\n);\n\nCREATE TABLE IF NOT EXISTS Messages (\n    message_id VARCHAR(1024) NOT NULL CONSTRAINT Messages_pk PRIMARY KEY,\n    content TEXT,\n    date INTEGER NOT NULL,\n    operation INTEGER NOT NULL,\n    sender_peer_id VARCHAR(1024) NOT NULL CONSTRAINT Messages_Peers_peer_id_fk REFERENCES Peers ON UPDATE CASCADE,\n    chat_id VARCHAR(1024) NOT NULL CONSTRAINT Messages_Chats_chat_id_fk REFERENCES Chats ON UPDATE CASCADE,\n    receiver_peer_id VARCHAR(1024) NOT NULL CONSTRAINT Messages_Peers_peer_id_fk_2 REFERENCES Peers,\n    sender_address VARCHAR(1024) NOT NULL,\n    receiver_address VARCHAR(1024) NOT NULL\n);\n\nCREATE TABLE IF NOT EXISTS Invitations (\n    invitation_id INTEGER NOT NULL CONSTRAINT Invitations_pk PRIMARY KEY AUTOINCREMENT,\n    invitation_status INTEGER NOT NULL,\n    message_id VARCHAR(1024) NOT NULL CONSTRAINT Invitations_Messages_message_id_fk REFERENCES Messages\n);\n\nCREATE TABLE IF NOT EXISTS PeersInInvitedChat (\n    public_key VARCHAR(1024) NOT NULL,\n    invited_peer_id INTEGER NOT NULL CONSTRAINT PeersInInvitedChat_pk PRIMARY KEY AUTOINCREMENT,\n    address VARCHAR(1024) NOT NULL,\n    invitation_id INTEGER NOT NULL CONSTRAINT PeersInInvitedChat_Invitations_invitation_id_fk REFERENCES Invitations\n);\n"
	_, err := a.db.Exec(sqlCommands)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *StorageSQLiteAdapter) ChatCreated(chatName string, chatId string) error {
	//TODO implement me
	panic("implement me")
}

func (a *StorageSQLiteAdapter) PeerSetUsername(peerId string, chatId string, username string) error {
	//TODO implement me
	panic("implement me")
}

func (a *StorageSQLiteAdapter) SetPeerUsername(username, peerID, chatID string) error {
	stmt, err := a.db.Prepare("UPDATE ChatMembers SET username = ? WHERE peer_id = ? AND chat_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, peerID, chatID)
	return err
}

func (a *StorageSQLiteAdapter) PeerJoinedChat(timestamp int64, peerID, chatID string) error {
	stmt, err := a.db.Prepare("INSERT INTO ChatMembers (date, peer_id, chat_id) SELECT ?, ?, ? WHERE NOT EXISTS (SELECT 1 FROM ChatMembers WHERE peer_id = ? AND chat_id = ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(timestamp, peerID, chatID, peerID, chatID)
	return err
}

func (a *StorageSQLiteAdapter) PeerLeftChat(peerID, chatID string) error {
	stmt, err := a.db.Prepare("DELETE FROM ChatMembers WHERE peer_id = ? AND chat_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(peerID, chatID)
	return err
}

func (a *StorageSQLiteAdapter) CreateChat(chatID, name string) error {
	stmt, err := a.db.Prepare("INSERT INTO Chats (chat_id, name) SELECT ?, ? WHERE NOT EXISTS (SELECT 1 FROM Chats WHERE chat_id = ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(chatID, name, chatID)
	return err
}

func (a *StorageSQLiteAdapter) InvitedToChat(messageID string, peers []store.PublicKeyAddress) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT INTO Invitations (invitation_status, message_id) SELECT 0, ? WHERE NOT EXISTS (SELECT 1 FROM Invitations WHERE message_id = ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(messageID, messageID)
	if err != nil {
		return err
	}

	stmt, err = tx.Prepare("INSERT INTO PeersInInvitedChat (public_key, address, invitation_id) SELECT ?, ?, (SELECT invitation_id FROM Invitations WHERE message_id = ?) WHERE NOT EXISTS (SELECT 1 FROM PeersInInvitedChat WHERE public_key = ? AND address = ? AND invitation_id = (SELECT invitation_id FROM Invitations WHERE message_id = ?))")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, peer := range peers {
		_, err = stmt.Exec(peer.PublicKey, peer.Address, messageID, peer.PublicKey, peer.Address, messageID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (a *StorageSQLiteAdapter) PeerGotInvitedToChat(peerId string, chatId string) error {
	// Implement the logic to store a peer getting invited to a chat in the database
	// ...
	return nil
}

// TODO: Rework
func (a *StorageSQLiteAdapter) GetInvitations(peerID string) ([]string, error) {
	rows, err := a.db.Query(`
		SELECT m.chat_id
		FROM Messages m
         JOIN Invitations i ON i.message_id = m.message_id
         JOIN Peers p ON m.receiver_peer_id = p.peer_id
		WHERE p.public_key = ? AND m.operation = ?
	`, peerID, network.INVITE_TO_CHAT)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var invitations []string
	for rows.Next() {
		var invitation string
		err := rows.Scan(&invitation)
		if err != nil {
			return nil, err
		}
		invitations = append(invitations, invitation)
	}

	return invitations, nil
}

// GetMissingInternalMessages retrieves the message IDs from inputMessageIDs that are not present in the Messages table for the specified chatID.
func (a *StorageSQLiteAdapter) GetMissingInternalMessages(chatID string, inputMessageIDs []string) ([]string, error) {
	if len(inputMessageIDs) == 0 {
		return nil, nil
	}

	// Prepare the query
	query := `SELECT message_id FROM Messages WHERE chat_id = ?`

	// Execute the query
	rows, err := a.db.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the results
	var existingMessageIDs []string
	for rows.Next() {
		var messageID string
		err := rows.Scan(&messageID)
		if err != nil {
			return nil, err
		}
		existingMessageIDs = append(existingMessageIDs, messageID)
	}

	diff := difference(inputMessageIDs, existingMessageIDs)

	return diff, nil
}

func difference(inputMessageIDs, existingMessageIDs []string) []string {
	idSet := make(map[string]struct{}, len(existingMessageIDs))
	for _, id := range existingMessageIDs {
		idSet[id] = struct{}{}
	}

	var diff []string
	for _, id := range inputMessageIDs {
		if _, found := idSet[id]; !found {
			diff = append(diff, id)
		}
	}

	return diff
}

func (a *StorageSQLiteAdapter) GetMissingExternalMessages(chatID string, inputMessageIDs []string) ([]string, error) {
	query := "SELECT message_id FROM Messages WHERE chat_id = ? AND message_id NOT IN (SELECT message_id FROM Messages WHERE chat_id = ? AND message_id IN (?" + strings.Repeat(",?", len(inputMessageIDs)-1) + "))"
	args := make([]interface{}, len(inputMessageIDs)+2)
	args[0] = chatID
	args[1] = chatID
	for i, id := range inputMessageIDs {
		args[i+2] = id
	}

	rows, err := a.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var message string
		err := rows.Scan(&message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (a *StorageSQLiteAdapter) StoreMessage(message network.Message) error {
	// Check if sender exists in Peers table, and insert if not
	err := a.insertPeerIfNotExists(message.SenderID, message.SenderAddress)
	if err != nil {
		return err
	}

	// Check if receiver exists in Peers table, and insert if not
	err = a.insertPeerIfNotExists(message.ReceiverID, message.ReceiverAddress)
	if err != nil {
		return err
	}

	stmt, err := a.db.Prepare(`
        INSERT INTO Messages (message_id, date, content, sender_peer_id, receiver_peer_id, sender_address, receiver_address, chat_id, operation)
        SELECT ?, ?, ?, (SELECT peer_id FROM Peers WHERE public_key = ?), (SELECT peer_id FROM Peers WHERE public_key = ?), ?, ?, ?, ?
        WHERE NOT EXISTS (SELECT 1 FROM Messages WHERE message_id = ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		message.Id, message.Timestamp, message.Content, message.SenderID, message.ReceiverID,
		message.SenderAddress, message.ReceiverAddress, message.ChatID, message.Operation, message.Id,
	)
	return err
}

func (a *StorageSQLiteAdapter) insertPeerIfNotExists(publicKey, address string) error {
	stmt, err := a.db.Prepare(`
        INSERT INTO Peers (public_key, address)
        SELECT ?, ?
        WHERE NOT EXISTS (SELECT 1 FROM Peers WHERE public_key = ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(publicKey, address, publicKey)
	return err
}

func (a *StorageSQLiteAdapter) RetrieveMessage(messageID string) (network.Message, error) {
	row := a.db.QueryRow(`
		SELECT m.message_id, m.date, m.content, m.operation, p.public_key, p2.public_key, m.sender_address, m.receiver_address, m.chat_id
		FROM Messages m, Peers p, Peers p2
		WHERE message_id = ? AND m.sender_peer_id = p.peer_id AND m.receiver_peer_id = p2.peer_id
	`, messageID)

	var message network.Message
	err := row.Scan(
		&message.Id, &message.Timestamp, &message.Content, &message.Operation, &message.SenderID, &message.ReceiverID,
		&message.SenderAddress, &message.ReceiverAddress, &message.ChatID,
	)
	if err != nil {
		return network.Message{}, err
	}

	return message, nil
}

func (a *StorageSQLiteAdapter) GetChats() ([]store.Chat, error) {
	rows, err := a.db.Query("SELECT chat_id, name FROM Chats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []store.Chat
	for rows.Next() {
		var chat store.Chat
		err := rows.Scan(&chat.ChatId, &chat.ChatName)
		if err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

func (a *StorageSQLiteAdapter) GetUsername(peerID, chatID string) (string, error) {
	row := a.db.QueryRow("SELECT username FROM ChatMembers WHERE peer_id = ? AND chat_id = ?", peerID, chatID)

	var username string
	err := row.Scan(&username)
	if err != nil {
		return "", err
	}

	return username, nil
}

func (a *StorageSQLiteAdapter) GetUsersInChat(chatID string) ([]store.User, error) {
	rows, err := a.db.Query("SELECT peer_id, username FROM ChatMembers WHERE chat_id = ?", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []store.User
	for rows.Next() {
		var user store.User
		var username sql.NullString
		err := rows.Scan(&user.UserId, &username)
		if err != nil {
			return nil, err
		}
		// Check if username is valid
		if username.Valid {
			user.Username = username.String
		} else {
			user.Username = "" // or any default value you prefer for NULL usernames
		}
		users = append(users, user)
	}

	return users, nil
}

func (a *StorageSQLiteAdapter) GetPeers() ([]string, error) {
	rows, err := a.db.Query("SELECT public_key FROM Peers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var peers []string
	for rows.Next() {
		var peer string
		err := rows.Scan(&peer)
		if err != nil {
			return nil, err
		}
		peers = append(peers, peer)
	}

	return peers, nil
}

func (a *StorageSQLiteAdapter) GetChatMessages(chatID string) ([]network.Message, error) {
	rows, err := a.db.Query(`
		SELECT m.message_id, m.content, m.date, m.operation, p.public_key, m.chat_id, p2.public_key, m.sender_address, m.receiver_address
		FROM Messages m, Peers p, Peers p2
		WHERE chat_id = ? AND m.sender_peer_id = p.peer_id AND m.receiver_peer_id = p2.peer_id
	`, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []network.Message
	for rows.Next() {
		var message network.Message
		err := rows.Scan(
			&message.Id, &message.Content, &message.Timestamp, &message.Operation, &message.SenderID, &message.ChatID, &message.ReceiverID,
			&message.SenderAddress, &message.ReceiverAddress,
		)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (a *StorageSQLiteAdapter) JoinChat(peerID, chatID string, invitationID int) error {
	// Check if the invitation exists and get the chat name
	var chatName string
	err := a.db.QueryRow("SELECT m.chat_id, c.name FROM Invitations i JOIN Messages m ON i.message_id = m.message_id JOIN Chats c ON m.chat_id = c.chat_id WHERE i.invitation_id = ? AND i.invitation_status = 0", invitationID).Scan(&chatID, &chatName)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("invitation not found or already accepted")
		}
		return err
	}

	// Check if the chat exists, and create it if not
	err = a.createChatIfNotExists(chatID, chatName)
	if err != nil {
		return err
	}

	// Check if the peer exists, and insert if not
	err = a.insertPeerIfNotExists(peerID, "")
	if err != nil {
		return err
	}

	// Insert the peer into the chat
	stmt, err := a.db.Prepare(`
        INSERT INTO ChatMembers (date, peer_id, chat_id)
        SELECT datetime('now'), (SELECT peer_id FROM Peers WHERE public_key = ?), ?
        WHERE NOT EXISTS (SELECT 1 FROM ChatMembers WHERE peer_id = (SELECT peer_id FROM Peers WHERE public_key = ?) AND chat_id = ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(peerID, chatID, peerID, chatID)
	if err != nil {
		return err
	}

	// Update the invitation status
	_, err = a.db.Exec("UPDATE Invitations SET invitation_status = 1 WHERE invitation_id = ?", invitationID)
	return err
}

func (a *StorageSQLiteAdapter) createChatIfNotExists(chatID, chatName string) error {
	stmt, err := a.db.Prepare(`
        INSERT INTO Chats (chat_id, name)
        SELECT ?, ?
        WHERE NOT EXISTS (SELECT 1 FROM Chats WHERE chat_id = ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(chatID, chatName, chatID)
	return err
}

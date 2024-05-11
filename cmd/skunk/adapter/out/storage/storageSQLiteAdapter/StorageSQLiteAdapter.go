package storageSQLiteAdapter

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/scherzma/Skunk/cmd/skunk/application/port/network"
	"log"
	"strings"
)

type StorageSQLiteAdapter struct {
	db *sql.DB
}

func NewStorageSQLiteAdapter(dbPath string) *StorageSQLiteAdapter {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	// Create tables if they don't exist
	adapter := &StorageSQLiteAdapter{db: db}
	adapter.CreateTables()

	return adapter
}

func (a *StorageSQLiteAdapter) CreateTables() {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS Chats (
			chat_id varchar(1024) not null constraint Chats_pk primary key,
			name varchar(40) not null
		)`,
		`CREATE TABLE IF NOT EXISTS Peers (
			peer_id integer not null constraint Peers_pk primary key autoincrement,
			public_key varchar(1024) not null,
			address varchar(1024) not null
		)`,
		`CREATE TABLE IF NOT EXISTS ChatMembers (
			chat_member_id integer not null constraint ChatMembers_pk primary key autoincrement,
			date datetime not null,
			peer_id varchar(1024) not null constraint ChatMembers_Peers_peer_id_fk references Peers on update cascade,
			chat_id varchar(1024) not null constraint ChatMembers_Chats_chat_id_fk references Chats on update cascade,
			username varchar(50)
		)`,
		`CREATE TABLE IF NOT EXISTS Messages (
			message_id varchar(1024) not null constraint Messages_pk primary key,
			content text,
			date datetime not null,
			operation integer not null,
			sender_peer_id varchar(1024) not null constraint Messages_Peers_peer_id_fk references Peers on update cascade,
			chat_id varchar(1024) not null constraint Messages_Chats_chat_id_fk references Chats on update cascade,
			receiver_peer_id varchar(1024) not null constraint Messages_Peers_peer_id_fk_2 references Peers,
			sender_address varchar(1024) not null,
			receiver_address varchar(1024) not null
		)`,
		`CREATE TABLE IF NOT EXISTS Invitations (
			invitation_id integer not null constraint Invitations_pk primary key autoincrement,
			invitation_status integer not null,
			message_id varchar(1024) not null constraint Invitations_Messages_message_id_fk references Messages
		)`,
		`CREATE TABLE IF NOT EXISTS PeersInInvitedChat (
			public_key varchar(1024) not null,
			invited_peer_id integer not null constraint PeersInInvitedChat_pk primary key autoincrement,
			address varchar(1024) not null,
			invitation_id integer not null constraint PeersInInvitedChat_Invitations_invitation_id_fk references Invitations
		)`,
	}

	for _, table := range tables {
		_, err := a.db.Exec(table)
		if err != nil {
			log.Fatal(err)
		}
	}
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

func (a *StorageSQLiteAdapter) PeerJoinedChat(peerID, chatID string) error {
	stmt, err := a.db.Prepare("INSERT INTO ChatMembers (date, peer_id, chat_id) SELECT datetime('now'), ?, ? WHERE NOT EXISTS (SELECT 1 FROM ChatMembers WHERE peer_id = ? AND chat_id = ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(peerID, chatID, peerID, chatID)
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

type PublicKeyAddress struct {
	Address   string
	PublicKey string
}

func (a *StorageSQLiteAdapter) InvitedToChat(messageID string, peers []PublicKeyAddress) error {
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

func (a *StorageSQLiteAdapter) PeerGotInvitedToChat(peerID, chatID string) error {
	// Implement the logic to store a peer getting invited to a chat in the database
	// ...
	return nil
}

func (a *StorageSQLiteAdapter) GetInvitations(peerID string) ([]string, error) {
	rows, err := a.db.Query(`
		SELECT m.*, i.invitation_status 
		FROM Messages m
		JOIN Invitations i ON i.message_id = m.message_id
		WHERE m.receiver_peer_id = ? AND m.operation = ?
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

func (a *StorageSQLiteAdapter) GetMissingInternalMessages(chatID string, inputMessageIDs []string) ([]string, error) {
	query := "SELECT message_id FROM Messages WHERE chat_id = ? AND message_id NOT IN (?" + strings.Repeat(",?", len(inputMessageIDs)-1) + ")"
	args := make([]interface{}, len(inputMessageIDs)+1)
	args[0] = chatID
	for i, id := range inputMessageIDs {
		args[i+1] = id
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
        SELECT ?, datetime('now'), ?, (SELECT peer_id FROM Peers WHERE public_key = ?), (SELECT peer_id FROM Peers WHERE public_key = ?), ?, ?, ?, ?
        WHERE NOT EXISTS (SELECT 1 FROM Messages WHERE message_id = ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		message.Id, message.Content, message.SenderID, message.ReceiverID,
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
	row := a.db.QueryRow("SELECT * FROM Messages WHERE message_id = ?", messageID)

	var message network.Message
	err := row.Scan(
		&message.Id, &message.Content, &message.Timestamp, &message.Operation, &message.SenderID, &message.ChatID, &message.ReceiverID,
		&message.SenderAddress, &message.ReceiverAddress,
	)
	if err != nil {
		return network.Message{}, err
	}

	return message, nil
}

type Chat struct {
	ChatID string
	Name   string
}

func (a *StorageSQLiteAdapter) GetChats() ([]Chat, error) {
	rows, err := a.db.Query("SELECT chat_id, name FROM Chats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.ChatID, &chat.Name)
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

type User struct {
	PeerID   string
	ChatID   string
	Username string
}

func (a *StorageSQLiteAdapter) GetUsersInChat(chatID string) ([]User, error) {
	rows, err := a.db.Query("SELECT peer_id, chat_id, username FROM ChatMembers WHERE chat_id = ?", chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.PeerID, &user.ChatID, &user.Username)
		if err != nil {
			return nil, err
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
		SELECT message_id, content, date, operation, sender_peer_id, chat_id, receiver_peer_id, sender_address, receiver_address
		FROM Messages
		WHERE chat_id = ?
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

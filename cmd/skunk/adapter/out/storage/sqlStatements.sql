
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

-- Peer Set Username
UPDATE ChatMembers
SET username = ?
WHERE peer_id = ? AND chat_id = ?;

-- Peer Joined Chat
INSERT INTO ChatMembers (peer_id, chat_id)
SELECT ?, ?
    WHERE NOT EXISTS(SELECT 1 FROM ChatMembers WHERE peer_id=? AND chat_id=?);

-- Peer Left Chat
DELETE FROM ChatMembers
WHERE peer_id=? AND chat_id=?;

-- Chat Created
INSERT INTO Chats (chat_id, name)
SELECT ?, ?
    WHERE NOT EXISTS(SELECT 1 FROM Chats WHERE chat_id=? AND name=?);

-- Invited To Chat
-- type publicKeyAddress struct {
-- 	Address string
-- 	PublicKey string
-- }
-- InvitatedToChat(messageId string, peers []publicKeyAddress) error
INSERT INTO Invitations (message_id)
SELECT ?
    WHERE NOT EXISTS(SELECT 1 FROM Invitations WHERE message_id = ?);

-- For each peer
INSERT INTO PeersInInvitedChat (public_key, address, invitation_id)
SELECT ?, ?, (SELECT Invitations.invitation_id FROM Invitations WHERE message_id = ?)
    WHERE NOT EXISTS(SELECT 1 FROM PeersInInvitedChat WHERE public_key = ? AND address = ? AND invitation_id = (SELECT Invitations.invitation_id FROM Invitations WHERE message_id = ?));

-- PeerGotInvitedToChat
-- PeerGotInvitedToChat(message network.Message) error


-- GetInvitations(peerId string) []string
SELECT m.*, i.invitation_status FROM Messages m, Invitations i WHERE m.receiver_peer_id = ?
                                                                 AND i.message_id = m.message_id AND m.operation = ?;

-- GetMissingInternalMessages(chatId string, inputMessageIDs []string) []string
SELECT * FROM Messages
WHERE message_id NOT IN (?, ?) AND chat_id = ?;

-- GetMissingExternalMessages(chatId string, inputMessageIDs []string) []string
SELECT * FROM Messages
WHERE (?, ?) NOT IN (SELECT message_id FROM Messages WHERE chat_id = ?) AND chat_id = ?;

-- StoreMessage(message network.Message) error
INSERT INTO Messages (message_id, date, content, sender_peer_id, receiver_peer_id, sender_address, receiver_address, chat_id, operation)
SELECT ?, ?, ?, ?, ?, ? ,?, ?, ?
    WHERE NOT EXISTS(SELECT 1 FROM Messages WHERE message_id = ?);

-- RetrieveMessage(messageId string) (network.Message, error)
SELECT * FROM Messages WHERE message_id = ?;

-- GetChats() []Chat
SELECT * FROM Chats;

-- GetUsername(peerId string, chatId string) string
SELECT username FROM ChatMembers WHERE peer_id = ? AND chat_id = ?;

-- GetUsersInChat(chatId string) []User
SELECT * FROM ChatMembers WHERE chat_id = ?;

-- GetPeers() []string
SELECT * FROM Peers;

-- GetChatMessages(chatId string) []ChatMessage
SELECT * FROM Messages WHERE chat_id = ?;

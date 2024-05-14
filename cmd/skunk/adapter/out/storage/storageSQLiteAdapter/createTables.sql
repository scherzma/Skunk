CREATE TABLE IF NOT EXISTS Chats (
    chat_id VARCHAR(1024) NOT NULL CONSTRAINT Chats_pk PRIMARY KEY,
    name VARCHAR(40) NOT NULL
);

CREATE TABLE IF NOT EXISTS Peers (
    peer_id INTEGER NOT NULL CONSTRAINT Peers_pk PRIMARY KEY AUTOINCREMENT,
    public_key VARCHAR(1024) NOT NULL,
    address VARCHAR(1024) NOT NULL
);

CREATE TABLE IF NOT EXISTS ChatMembers (
    chat_member_id INTEGER NOT NULL CONSTRAINT ChatMembers_pk PRIMARY KEY AUTOINCREMENT,
    date INTEGER NOT NULL,
    peer_id VARCHAR(1024) NOT NULL CONSTRAINT ChatMembers_Peers_peer_id_fk REFERENCES Peers ON UPDATE CASCADE,
    chat_id VARCHAR(1024) NOT NULL CONSTRAINT ChatMembers_Chats_chat_id_fk REFERENCES Chats ON UPDATE CASCADE,
    username VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS Messages (
    message_id VARCHAR(1024) NOT NULL CONSTRAINT Messages_pk PRIMARY KEY,
    content TEXT,
    date INTEGER NOT NULL,
    operation INTEGER NOT NULL,
    sender_peer_id VARCHAR(1024) NOT NULL CONSTRAINT Messages_Peers_peer_id_fk REFERENCES Peers ON UPDATE CASCADE,
    chat_id VARCHAR(1024) NOT NULL CONSTRAINT Messages_Chats_chat_id_fk REFERENCES Chats ON UPDATE CASCADE,
    receiver_peer_id VARCHAR(1024) NOT NULL CONSTRAINT Messages_Peers_peer_id_fk_2 REFERENCES Peers,
    sender_address VARCHAR(1024) NOT NULL,
    receiver_address VARCHAR(1024) NOT NULL
);

CREATE TABLE IF NOT EXISTS Invitations (
    invitation_id INTEGER NOT NULL CONSTRAINT Invitations_pk PRIMARY KEY AUTOINCREMENT,
    invitation_status INTEGER NOT NULL,
    message_id VARCHAR(1024) NOT NULL CONSTRAINT Invitations_Messages_message_id_fk REFERENCES Messages
);

CREATE TABLE IF NOT EXISTS PeersInInvitedChat (
    public_key VARCHAR(1024) NOT NULL,
    invited_peer_id INTEGER NOT NULL CONSTRAINT PeersInInvitedChat_pk PRIMARY KEY AUTOINCREMENT,
    address VARCHAR(1024) NOT NULL,
    invitation_id INTEGER NOT NULL CONSTRAINT PeersInInvitedChat_Invitations_invitation_id_fk REFERENCES Invitations
);

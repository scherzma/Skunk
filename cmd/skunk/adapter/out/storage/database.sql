create table Chats
(
    chat_id varchar(1024) not null
        constraint Chats_pk
            primary key,
    name    varchar(40)   not null
);

create table Peers
(
    peer_id    integer       not null
        constraint Peers_pk
            primary key autoincrement,
    public_key varchar(1024) not null,
    address    varchar(1024) not null
);

create table ChatMembers
(
    chat_member_id integer       not null
        constraint ChatMembers_pk
            primary key autoincrement,
    timestamp      integer       not null,
    peer_id        varchar(1024) not null
        constraint ChatMembers_Peers_peer_id_fk
            references Peers
            on update cascade,
    chat_id        varchar(1024) not null
        constraint ChatMembers_Chats_chat_id_fk
            references Chats
            on update cascade,
    username       varchar(50)
);

create table Messages
(
    message_id       varchar(1024) not null
        constraint Messages_pk
            primary key,
    content          text,
    timestamp        integer       not null,
    operation        integer       not null,
    sender_peer_id   varchar(1024) not null
        constraint Messages_Peers_peer_id_fk
            references Peers
            on update cascade,
    chat_id          varchar(1024) not null
        constraint Messages_Chats_chat_id_fk
            references Chats
            on update cascade,
    receiver_peer_id varchar(1024) not null
        constraint Messages_Peers_peer_id_fk_2
            references Peers,
    sender_address   varchar(1024) not null,
    receiver_address varchar(1024) not null
);

create table Invitations
(
    invitation_id     integer       not null
        constraint Invitations_pk
            primary key autoincrement,
    invitation_status integer       not null,
    message_id        varchar(1024) not null
        constraint Invitations_Messages_message_id_fk
            references Messages
);

create table PeersInInvitedChat
(
    public_key      varchar(1024) not null,
    invited_peer_id integer       not null
        constraint PeersInInvitedChat_pk
            primary key autoincrement,
    address         varchar(1024) not null,
    invitation_id   integer       not null
        constraint PeersInInvitedChat_Invitations_invitation_id_fk
            references Invitations
);


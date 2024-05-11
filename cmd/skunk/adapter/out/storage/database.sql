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
    date           datetime      not null,
    peer_id        varchar(1024) not null
        constraint ChatMembers_Peers_peer_id_fk
            references Peers
            on update cascade,
    chat_id        varchar(1024) not null
        constraint ChatMembers_Chats_chat_id_fk
            references Chats
            on update cascade
);

create table Invitations
(
    invitation_id     integer       not null
        constraint Invitations_pk
            primary key autoincrement,
    sender_peer_id    varchar(1024) not null
        constraint Invitations_Peers_peer_id_fk
            references Peers
            on update cascade,
    recipient_peer_id varchar(1024) not null
        constraint Invitations_Peers_peer_id_fk_2
            references Peers
            on update cascade,
    chat_id           varchar(1024) not null
        constraint Invitations_Chats_chat_id_fk
            references Chats
            on update cascade,
    date              datetime      not null,
    invitation_status integer       not null
);

create table Messages
(
    message_id     integer       not null
        constraint Messages_pk
            primary key autoincrement,
    content        text,
    date           datetime      not null,
    operation      integer       not null,
    sender_peer_id varchar(1024) not null
        constraint Messages_Peers_peer_id_fk
            references Peers
            on update cascade,
    chat_id        varchar(1024) not null
        constraint Messages_Chats_chat_id_fk
            references Chats
            on update cascade
);

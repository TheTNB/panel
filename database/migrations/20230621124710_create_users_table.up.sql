CREATE TABLE users
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    username   varchar(255)                      NOT NULL,
    password   varchar(255)                      NOT NULL,
    email      varchar(255) DEFAULT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

CREATE UNIQUE INDEX users_username_unique ON users (username);

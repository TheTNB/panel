CREATE TABLE cert_dns
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    type       varchar(255)                      NOT NULL,
    data       text                              NOT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

CREATE TABLE databases
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    name       varchar(255)                      NOT NULL,
    type       varchar(255)                      NOT NULL,
    host       varchar(255)                      NOT NULL,
    port       integer                           NOT NULL,
    username   varchar(255)                      NOT NULL,
    password   varchar(255) DEFAULT NULL,
    remark     text         DEFAULT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

CREATE UNIQUE INDEX databases_name_unique ON databases (name);
CREATE INDEX databases_type_index ON databases (type);

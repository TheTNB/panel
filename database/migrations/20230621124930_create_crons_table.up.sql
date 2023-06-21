CREATE TABLE crons
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    name       varchar(255)                      NOT NULL,
    status     boolean      DEFAULT 0            NOT NULL,
    type       varchar(255)                      NOT NULL,
    time       varchar(255)                      NOT NULL,
    shell      varchar(255) DEFAULT NULL,
    log        varchar(255) DEFAULT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

CREATE UNIQUE INDEX crons_name_unique ON crons (name);
CREATE INDEX crons_status_index ON crons (status);
CREATE INDEX crons_type_index ON crons (type);

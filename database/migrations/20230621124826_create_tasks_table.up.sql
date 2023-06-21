CREATE TABLE tasks
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    name       varchar(255)                      NOT NULL,
    status     varchar(255) DEFAULT 'waiting'    NOT NULL,
    shell      varchar(255) DEFAULT NULL,
    log        varchar(255) DEFAULT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

CREATE INDEX tasks_status_index ON tasks (status);

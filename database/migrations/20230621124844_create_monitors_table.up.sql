CREATE TABLE monitors
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    info       text                              NOT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

CREATE TABLE settings
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    key        varchar(255)                      NOT NULL,
    value      varchar(255) DEFAULT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

CREATE UNIQUE INDEX settings_key_unique ON settings (key);

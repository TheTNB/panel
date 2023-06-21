CREATE TABLE websites
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    name       varchar(255)                      NOT NULL,
    status     boolean DEFAULT 1                 NOT NULL,
    path       varchar(255)                      NOT NULL,
    php        integer DEFAULT 0                 NOT NULL,
    ssl        boolean DEFAULT 0                 NOT NULL,
    remark     text    DEFAULT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

CREATE UNIQUE INDEX websites_name_unique ON websites (name);
CREATE INDEX websites_status_index ON websites (status);
CREATE INDEX websites_php_index ON websites (php);
CREATE INDEX websites_ssl_index ON websites (ssl);

CREATE TABLE certs
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    user_id    integer                           NOT NULL,
    website_id integer      DEFAULT NULL,
    dns_id     integer      DEFAULT NULL,
    cron_id    integer      DEFAULT NULL,
    type       varchar(255)                      NOT NULL,
    domains    text                              NOT NULL,
    cert_url   varchar(255) DEFAULT NULL,
    cert       text         DEFAULT NULL,
    key        text         DEFAULT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

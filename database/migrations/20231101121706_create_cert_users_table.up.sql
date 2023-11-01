CREATE TABLE cert_users
(
    id           integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    email        varchar(255)                      NOT NULL,
    ca           varchar(255)                      NOT NULL,
    kid          varchar(255) DEFAULT NULL,
    hmac_encoded varchar(255) DEFAULT NULL,
    private_key  text                              NOT NULL,
    key_type     varchar(255)                      NOT NULL,
    created_at   datetime                          NOT NULL,
    updated_at   datetime                          NOT NULL
);

CREATE INDEX idx_cert_users_email ON cert_users (email);

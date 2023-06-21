CREATE TABLE plugins
(
    id         integer PRIMARY KEY AUTOINCREMENT NOT NULL,
    slug       varchar(255)                      NOT NULL,
    version    varchar(255)                      NOT NULL,
    show       boolean DEFAULT 0                 NOT NULL,
    show_order integer DEFAULT 0                 NOT NULL,
    created_at datetime                          NOT NULL,
    updated_at datetime                          NOT NULL
);

CREATE UNIQUE INDEX plugins_slug_unique ON plugins (slug);
CREATE INDEX plugins_show_index ON plugins (show);
CREATE INDEX plugins_show_order_index ON plugins (show_order);

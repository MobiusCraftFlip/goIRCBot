-- Initial schema

CREATE TABLE bots (
    id         SERIAL      PRIMARY KEY,
    nick       TEXT        NOT NULL,
    username   TEXT        NOT NULL,
    realname   TEXT        NOT NULL,
    server     TEXT        NOT NULL,
    port       INTEGER     NOT NULL DEFAULT 6697,
    ssl        BOOLEAN     NOT NULL DEFAULT True,
    admin_chan TEXT        NOT NULL DEFAULT '#botteam',
    log_chan   TEXT        NOT NULL DEFAULT '#bots.log',
    modules    TEXT[]      NOT NULL DEFAULT '{}'
);

CREATE TABLE bot_channels (
    id     SERIAL  PRIMARY KEY,
    bot_id INTEGER NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    channel TEXT   NOT NULL,
    UNIQUE (bot_id, channel)
);

CREATE TABLE bot_config (
    id     SERIAL  PRIMARY KEY,
    bot_id INTEGER NOT NULL REFERENCES bots(id) ON DELETE CASCADE,
    key    TEXT    NOT NULL,
    value  TEXT    NOT NULL DEFAULT '',
    UNIQUE (bot_id, key)
);

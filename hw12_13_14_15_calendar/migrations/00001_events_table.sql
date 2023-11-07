-- +goose Up
CREATE TABLE events
(
    id          text,
    title       text,
    date        timestamptz,
    duration    int,
    description text,
    user_id     text,
    reminder    text
);

CREATE INDEX id_idx ON events (id);

-- +goose Down
DROP TABLE events;
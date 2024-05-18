-- +goose Up

ALTER TABLE sessions ADD column created_at timestamptz;

CREATE TABLE absence_reminder_messages
(
    user_id BIGINT PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE  NOT NULL,
    created_at timestamp NOT NULL
);


-- +goose Down
ALTER TABLE sessions DROP column created_at;

DROP TABLE absence_reminder_messages;
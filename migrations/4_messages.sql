-- +goose Up

CREATE TABLE messages
(
    user_id BIGINT  REFERENCES users (id) ON DELETE CASCADE  NOT NULL,
    message_type TEXT NOT NULL,
    created_at timestamp NOT NULL,
    PRIMARY KEY (user_id, message_type)
);

ALTER TABLE sessions ADD column created_at timestamptz;


-- +goose Down
DROP TABLE messages;

ALTER TABLE sessions DROP column created_at;


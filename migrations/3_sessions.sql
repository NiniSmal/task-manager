-- +goose Up
CREATE TABLE sessions
(
    id      uuid PRIMARY KEY,
    user_id BIGINT NOT NULL,
    role TEXT NOT NULL
);

-- +goose Down
DROP TABLE sessions;



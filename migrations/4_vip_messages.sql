-- +goose Up

CREATE TABLE vip_messages
(
    user_id BIGINT PRIMARY KEY REFERENCES users (id)ON DELETE CASCADE  NOT NULL,
    created_at timestamp NOT NULL
);


-- +goose Down
DROP TABLE vip_messages
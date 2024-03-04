-- +goose Up
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    login      TEXT        NOT NULL,
    password   TEXT        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE users;



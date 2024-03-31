-- +goose Up
CREATE TABLE users
(
    id         BIGSERIAL PRIMARY KEY,
    login      TEXT        NOT NULL UNIQUE,
    password   TEXT        NOT NULL,
    created_at timestamptz NOT NULL,
    role       TEXT        NOT NULL,
    verification BOOLEAN,
    verification_code TEXT
);

-- +goose Down
DROP TABLE users;



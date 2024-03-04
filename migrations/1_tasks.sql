-- +goose Up
CREATE TABLE tasks
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    status     TEXT        NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE tasks;
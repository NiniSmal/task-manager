--+goose Up
CREATE TABLE projects
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT        NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz NOT NULL,
    user_id    BIGINT
);

-- +goose Down
DROP TABLE projects;

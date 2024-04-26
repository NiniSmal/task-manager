-- +goose Up
CREATE TABLE tasks
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT                            NOT NULL,
    status     TEXT                            NOT NULL,
    created_at TIMESTAMPTZ                     NOT NULL,
    user_id    BIGINT REFERENCES users (id)    NOT NULL,
    project_id BIGINT REFERENCES projects (id) NOT NULL
);

-- +goose Down
DROP TABLE tasks;
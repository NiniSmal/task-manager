-- +goose Up
CREATE TABLE users
(
    id                BIGSERIAL PRIMARY KEY,
    email             TEXT        NOT NULL UNIQUE,
    password          TEXT        NOT NULL,
    created_at        timestamptz NOT NULL,
    role              TEXT        NOT NULL,
    verification      BOOLEAN     NOT NULL,
    verification_code TEXT        NOT NULL
);

CREATE TABLE sessions
(
    id      uuid PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) NOT NULL
);

CREATE TABLE projects
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT                         NOT NULL,
    created_at timestamptz                  NOT NULL,
    updated_at timestamptz                  NOT NULL,
    user_id    BIGINT REFERENCES users (id) NOT NULL
);

CREATE TABLE tasks
(
    id         BIGSERIAL PRIMARY KEY,
    name       TEXT                            NOT NULL,
    status     TEXT                            NOT NULL,
    created_at TIMESTAMPTZ                     NOT NULL,
    user_id    BIGINT REFERENCES users (id)    NOT NULL,
    project_id BIGINT REFERENCES projects (id) NOT NULL
);

CREATE TABLE user_projects
(
    user_id    BIGINT REFERENCES users (id)    NOT NULL,
    project_id BIGINT REFERENCES projects (id) NOT NULL,
    PRIMARY KEY (user_id, project_id)


);

-- +goose Down
DROP TABLE user_projects;
DROP TABLE tasks;
DROP TABLE projects;
DROP TABLE sessions;
DROP TABLE users;



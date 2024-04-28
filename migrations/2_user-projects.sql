-- +goose Up
CREATE TABLE user_projects
(
    project_id BIGINT REFERENCES projects (id) NOT NULL,
    user_id    BIGINT REFERENCES users (id)    NOT NULL


);

-- -goose Down
DROP TABLE user_projects;
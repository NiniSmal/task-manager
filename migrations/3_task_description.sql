-- +goose Up
ALTER TABLE tasks ADD COLUMN description text;

CREATE TABLE codes_projects_users
(
    code uuid PRIMARY KEY,
    project_id     BIGINT REFERENCES projects(id) ON DELETE CASCADE NOT NULL,
    user_id BIGINT REFERENCES users (id)ON DELETE CASCADE NOT NULL
);

-- +goose Down
ALTER TABLE tasks DROP COLUMN description;

DROP TABLE codes_projects_users;
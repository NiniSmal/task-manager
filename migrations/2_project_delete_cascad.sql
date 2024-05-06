-- +goose Up
ALTER TABLE projects DROP CONSTRAINT projects_user_id_fkey;
ALTER TABLE projects ADD CONSTRAINT projects_user_id_fkey foreign key(user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE tasks DROP CONSTRAINT tasks_project_id_fkey;
ALTER TABLE tasks ADD CONSTRAINT tasks_project_id_fkey foreign key(project_id) REFERENCES projects(id) ON DELETE CASCADE;

ALTER TABLE tasks DROP CONSTRAINT tasks_user_id_fkey;
ALTER TABLE tasks ADD CONSTRAINT tasks_user_id_fkey foreign key(user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE user_projects DROP CONSTRAINT user_projects_project_id_fkey;
ALTER TABLE user_projects ADD CONSTRAINT user_projects_project_id_fkey foreign key(project_id) REFERENCES projects(id) ON DELETE CASCADE;

ALTER TABLE user_projects DROP CONSTRAINT user_projects_user_id_fkey;
ALTER TABLE user_projects ADD CONSTRAINT user_projects_user_id_fkey foreign key(user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE sessions DROP CONSTRAINT sessions_user_id_fkey;
ALTER TABLE sessions ADD CONSTRAINT sessions_user_id_fkey foreign key(user_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down

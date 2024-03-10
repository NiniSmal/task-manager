-- +goose Up
ALTER TABLE tasks ADD COLUMN user_id BIGINT, ADD COLUMN role TEXT;

-- +goose Down
ALTER TABLE tasks DROP COLUMN user_id, DROP COLUMN role ;
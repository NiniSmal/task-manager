-- +goose Up
ALTER TABLE tasks ADD COLUMN user_id BIGINT;

-- +goose Down
ALTER TABLE tasks DROP COLUMN user_id;
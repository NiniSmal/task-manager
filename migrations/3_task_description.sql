-- +goose Up
ALTER TABLE tasks ADD COLUMN description text;


-- +goose Down
ALTER TABLE tasks DROP COLUMN description;
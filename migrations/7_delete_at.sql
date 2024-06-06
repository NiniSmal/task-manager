-- +goose Up

ALTER TABLE projects ADD COLUMN deleted_at timestamp;
ALTER TABLE tasks ADD COLUMN deleted_at timestamp;
ALTER TABLE users ADD COLUMN deleted_at timestamp;


-- +goose Down
ALTER TABLE projects DROP column deleted_at;
ALTER TABLE tasks DROP column deleted_at;
ALTER TABLE users DROP column deleted_at;
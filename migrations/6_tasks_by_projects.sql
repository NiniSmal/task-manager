-- +goose Up
    ALTER TABLE tasks ADD COLUMN  project_id BIGINT NOT NULL default 0;

--+goose Down
    ALTER TABLE tasks DROP COLUMN project_id;












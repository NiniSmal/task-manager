-- +goose Up
ALTER TABLE tasks ADD column assigner_id BIGINT NOT NULL DEFAULT 0;
UPDATE tasks SET assigner_id = user_id;
ALTER TABLE tasks ALTER column assigner_id DROP DEFAULT;

-- +goose Down
ALTER TABLE tasks DROP column assigner_id;
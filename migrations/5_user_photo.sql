-- +goose Up
ALTER TABLE users ADD column photo text;




-- +goose Down
ALTER TABLE users DROP column photo;
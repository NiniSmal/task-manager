-- +goose Up
ALTER TABLE users DROP column photo;

CREATE TABLE profiles
(
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE NOT NULL,
    image_url TEXT NOT NULL
);


-- +goose Down

DROP TABLE profiles

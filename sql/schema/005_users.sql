-- +goose Up
ALTER TABLE users
ADD is_chirpy_red Boolean NOT NULL DEFAULT 'false';

-- +goose Down
ALTER TABLE users
DROP COLUMN is_chirpy_red;
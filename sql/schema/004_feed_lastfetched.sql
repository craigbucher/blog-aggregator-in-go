-- Create a new migration that adds a last_fetched_at column to the feeds table. It should be nullable
-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP COLUMN last_fetched_at;
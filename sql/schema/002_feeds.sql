-- +goose Up
-- Create a feeds table. Like any table in our DB, we'll need the standard id, created_at, and 
-- updated_at fields; We'll also need a few more:
CREATE TABLE feeds (
    id UUID PRIMARY KEY,            
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    -- The name of the feed (like "The Changelog, or "The Boot.dev Blog")
    name TEXT NOT NULL,
    -- The URL of the feed
    -- Make the url field unique so that in the future we aren't downloading duplicate posts
    url TEXT NOT NULL UNIQUE,
    -- The ID of the user who added this feed
    -- Use an ON DELETE CASCADE constraint on the user_id foreign key so that if a user is deleted, 
    -- all of their feeds are automatically deleted as well
    -- So, deleting the users in the reset command also deletes all of their feeds
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
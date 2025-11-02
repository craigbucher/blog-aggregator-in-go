-- Add a posts table to the database:
-- A post is a single entry from a feed. It should have:
-- +goose Up
CREATE TABLE posts (
    id UUID PRIMARY KEY,    -- a unique identifier for the post
    created_at TIMESTAMP NOT NULL,  -- the time the record was created
    updated_at TIMESTAMP NOT NULL,  -- the time the record was last updated
    title TEXT NOT NULL,    -- the title of the post
    url TEXT NOT NULL UNIQUE,   -- the URL of the post (this should be unique)
    description TEXT,   -- the description of the post
    published_at TIMESTAMP, -- the time the post was published
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE    -- the ID of the feed that the post came from
);

-- +goose Down
DROP TABLE posts;
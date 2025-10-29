-- +goose Up
-- Create a feed_follows table with a new migration. It should:
CREATE TABLE feed_follows (
    -- Have an id column that is a primary key:
    id UUID PRIMARY KEY,
    -- Have created_at and updated_at columns
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    -- Have user_id and feed_id foreign key columns. Feed follows should auto delete when a 
    -- user or feed is deleted:
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    -- Add a unique constraint on user/feed pairs - we don't want duplicate follow records:
    UNIQUE (user_id, feed_id)
);

-- +goose Down
DROP TABLE feed_follows;
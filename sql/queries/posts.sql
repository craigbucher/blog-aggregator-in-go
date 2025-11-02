-- Add a "create post" SQL query to the database. This should insert 
-- a new post into the database:
-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;
--
-- Add a "get posts for user" SQL query to the database:
-- name: GetPostsForUser :many
SELECT posts.*, feeds.name AS feed_name FROM posts
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
JOIN feeds ON posts.feed_id = feeds.id
WHERE feed_follows.user_id = $1
-- Order the results so that the most recent posts are first:
ORDER BY posts.published_at DESC
-- Make the number of posts returned configurable:
LIMIT $2;
--
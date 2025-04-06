-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetUserPosts :many
SELECT posts.*, feeds.name as feed_name FROM users
INNER JOIN feed_follows ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
INNER JOIN posts ON feeds.id = posts.id
WHERE users.id = $1
ORDER BY posts.published_at DESC
LIMIT $2;
-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT 
    inserted_feed_follow.*, 
    feeds.name AS feed_name, 
    users.name AS user_name
FROM inserted_feed_follow
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
    (
        SELECT name 
        FROM users
        WHERE id = user_id
    ) as user_name,
    (
        SELECT name
        FROM feeds
        WHERE id = feed_id
    ) as feed_name
FROM feed_follows
WHERE feed_follows.user_id = $1;

-- name: DeteleFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_id = $1 AND user_id = $2;
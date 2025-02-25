-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeedsWithUsers :many
SELECT feeds.name AS feed_name, feeds.url, users.name AS user_name
FROM feeds
JOIN users ON feeds.user_id = users.id;


-- name: GetFeedByURL :one
SELECT * FROM feeds
WHERE url = $1
LIMIT 1;


-- name: CreateFeedFollow :one
WITH inserted_follow AS (
  INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
  VALUES ($1, $2, $3, $4, $5)
  RETURNING *
)
SELECT 
  inserted_follow.*,
  users.name AS user_name,
  feeds.name AS feed_name
FROM inserted_follow
JOIN users ON inserted_follow.user_id = users.id
JOIN feeds ON inserted_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT 
  feed_follows.*,
  users.name AS user_name,
  feeds.name AS feed_name
FROM feed_follows
JOIN users ON feed_follows.user_id = users.id
JOIN feeds ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY feed_follows.created_at DESC;


-- name: DeleteFeedFollowByUserAndFeedURL :exec
DELETE FROM feed_follows
USING feeds
WHERE feed_follows.feed_id = feeds.id
AND feed_follows.user_id = $1
AND feeds.url = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;

-- name: DeleteAllFeeds :exec
DELETE FROM feeds;
-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;    

-- name: GetUser :one
SELECT * FROM users WHERE name = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetPostsForUser :many
SELECT posts.*
FROM posts
JOIN feed_follows ON posts.feed_id = feed_follows.feed_id
JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1
ORDER BY posts.published_at DESC
LIMIT $2;
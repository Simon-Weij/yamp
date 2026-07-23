-- name: CreateUser :one
INSERT INTO users (
    username, password_hash
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetUserForLogin :one
SELECT id, password_hash FROM users WHERE username = $1;

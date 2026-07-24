-- name: CreateUser :one
INSERT INTO users (
    username, password_hash
) VALUES (
    $1, $2
)
RETURNING *;

-- name: GetUserForLogin :one
SELECT id, password_hash FROM users WHERE username = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
    user_id, token_hash, expires_at
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetRefreshTokenByHash :one
SELECT id, user_id, token_hash, expires_at, revoked_at, created_at
FROM refresh_tokens
WHERE token_hash = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE id = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(created_at, updated_at, expires_at, revoked_at, token, user_id)
VALUES (NOW(), NOW(), NOW() + INTERVAL '60 days', NULL, $1, $2)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1;
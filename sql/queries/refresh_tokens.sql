-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days'
)
RETURNING *;

-- name: LookupRefreshToken :one
SELECT * FROM refresh_tokens
WHERE token = $1;

-- name: GetRefreshTokenFromUser :one
SELECT * FROM refresh_tokens
WHERE user_id = $1;

-- name: GetUserFromRefreshToken :one
SELECT * FROM users
INNER JOIN refresh_tokens
ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1;

-- name: SetRefreshTokenAsRevoked :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1;
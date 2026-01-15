-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, expires_at, revoked_at, user_id)
VALUES (
    $1,
    NOW() AT TIME ZONE 'UTC',
    NOW() AT TIME ZONE 'UTC',
    NOW() AT TIME ZONE 'UTC' + INTERVAL '60 days',
    NULL,
    $2
)
RETURNING token;


-- name: RefreshTokenLookup :one
SELECT expires_at, revoked_at, user_id FROM refresh_tokens
WHERE token = $1;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW() AT TIME ZONE 'UTC', updated_at = NOW() AT TIME ZONE 'UTC'
WHERE token = $1;
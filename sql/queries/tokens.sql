-- name: CreateRefreshToken :one 
INSERT INTO refresh_tokens(token, user_id, expires_at, revoked_at, created_at, updated_at)
VALUES (
  $1, 
  $2,
  $3, 
  NULL,
  NOW(),
  NOW()
)
RETURNING *;

-- name: GetRefreshTokenByToken :one 
SELECT token 
FROM refresh_tokens 
WHERE token = $1;

-- name: GetUserByRefreshToken :one 
SELECT users.* FROM users 
INNER JOIN refresh_tokens 
ON refresh_tokens.user_id = users.id
WHERE refresh_tokens.token = $1
AND refresh_tokens.revoked_at IS NULL
AND refresh_tokens.expires_at > NOW();

-- name: RevokeRefreshToken :exec 
UPDATE refresh_tokens 
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1; 


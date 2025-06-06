-- name: CreateSession :one
INSERT INTO sessions (
  id, user_id, refresh_token, user_agent, client_ip, is_blocked, expires_at
)
VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions WHERE id = $1;

-- name: BlockSessionByID :exec
UPDATE sessions 
SET is_blocked = true
WHERE id = $1;

-- name: GetReusableSessionByUserID :one
SELECT * FROM sessions
WHERE user_id = $1 AND is_blocked = false AND expires_at > $2
ORDER BY expires_at DESC
LIMIT 1;

-- name: UpdateSessionTokenAndExpiry :exec
UPDATE sessions
SET refresh_token = $1, expires_at = $2
WHERE id = $3;


-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = $1;

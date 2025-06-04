-- name: CreateMediaGroup :one
INSERT INTO media_groups (user_id, name)
VALUES ($1, $2)
RETURNING *;

-- name: GetGroupByID :one
SELECT * FROM media_groups
WHERE id = $1;

-- name: ListGroupsByUser :many
SELECT * FROM media_groups
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: DeleteMediaGroup :exec
DELETE FROM media_groups
WHERE id = $1;
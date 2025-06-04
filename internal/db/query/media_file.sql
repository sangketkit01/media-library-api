-- name: CreateMediaFile :one
INSERT INTO media_files (user_id, group_id, filename, file_type, size)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetMediaFileByID :one
SELECT * FROM media_files
WHERE id = $1;

-- name: ListMediaByUser :many
SELECT * FROM media_files
WHERE user_id = $1
ORDER BY uploaded_at DESC;

-- name: ListMediaByGroup :many
SELECT * FROM media_files
WHERE user_id = $1 AND group_id = $2
ORDER BY uploaded_at DESC;

-- name: AssignMediaToGroup :exec
UPDATE media_files
SET group_id = $2
WHERE id = $1;

-- name: CountMediaSizeByUser :one
SELECT COALESCE(SUM(size), 0) AS total_size
FROM media_files
WHERE user_id = $1;

-- name: DeleteMediaFile :exec
DELETE FROM media_files
WHERE id = $1;


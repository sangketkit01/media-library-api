// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: media_group.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createMediaGroup = `-- name: CreateMediaGroup :one
INSERT INTO media_groups (user_id, name)
VALUES ($1, $2)
RETURNING id, user_id, name, created_at
`

type CreateMediaGroupParams struct {
	UserID pgtype.UUID `json:"user_id"`
	Name   string      `json:"name"`
}

func (q *Queries) CreateMediaGroup(ctx context.Context, arg CreateMediaGroupParams) (MediaGroup, error) {
	row := q.db.QueryRow(ctx, createMediaGroup, arg.UserID, arg.Name)
	var i MediaGroup
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CreatedAt,
	)
	return i, err
}

const deleteMediaGroup = `-- name: DeleteMediaGroup :exec
DELETE FROM media_groups
WHERE id = $1
`

func (q *Queries) DeleteMediaGroup(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteMediaGroup, id)
	return err
}

const getGroupByID = `-- name: GetGroupByID :one
SELECT id, user_id, name, created_at FROM media_groups
WHERE id = $1
`

func (q *Queries) GetGroupByID(ctx context.Context, id pgtype.UUID) (MediaGroup, error) {
	row := q.db.QueryRow(ctx, getGroupByID, id)
	var i MediaGroup
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Name,
		&i.CreatedAt,
	)
	return i, err
}

const listGroupsByUser = `-- name: ListGroupsByUser :many
SELECT id, user_id, name, created_at FROM media_groups
WHERE user_id = $1
ORDER BY created_at DESC
`

func (q *Queries) ListGroupsByUser(ctx context.Context, userID pgtype.UUID) ([]MediaGroup, error) {
	rows, err := q.db.Query(ctx, listGroupsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []MediaGroup{}
	for rows.Next() {
		var i MediaGroup
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

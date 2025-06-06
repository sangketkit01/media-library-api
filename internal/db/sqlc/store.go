package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Queryer interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row 
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type Store interface {
	Querier
}

type SQLStore struct {
	*Queries
	db Queryer
}

func NewStore(db Queryer) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

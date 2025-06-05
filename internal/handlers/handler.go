package handlers

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sangketkit01/media-library-api/internal/config"
	db "github.com/sangketkit01/media-library-api/internal/db/sqlc"
	"github.com/sangketkit01/media-library-api/internal/token"
)

type Handler struct {
	Store  db.Store
	Config *config.Config
	tokenMaker token.Maker
}

func NewHandler(config *config.Config, tokenMaker token.Maker) (*Handler, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, config.DatabaseUrl)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	store := db.NewStore(conn)

	return &Handler{
		Config: config,
		Store:  store,
		tokenMaker: tokenMaker,
	}, nil
}

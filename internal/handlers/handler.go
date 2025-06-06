package handlers

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sangketkit01/media-library-api/internal/config"
	db "github.com/sangketkit01/media-library-api/internal/db/sqlc"
	"github.com/sangketkit01/media-library-api/internal/token"
)

type Handler struct {
	Store      db.Store
	Config     *config.Config
	tokenMaker token.Maker
	Pool *pgxpool.Pool
}

func NewHandler(config *config.Config, tokenMaker token.Maker) (*Handler, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cfg, err := pgxpool.ParseConfig(config.DatabaseUrl)
	if err != nil {
		return nil, err
	}

	cfg.MaxConns = 10
	cfg.MaxConnIdleTime = 5 * time.Minute
	cfg.MaxConnLifetime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil{
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	store := db.NewStore(pool)

	return &Handler{
		Config:     config,
		Store:      store,
		tokenMaker: tokenMaker,
		Pool: pool,
	}, nil
}

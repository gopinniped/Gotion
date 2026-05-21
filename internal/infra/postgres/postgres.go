package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/gopinniped/gotion/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	DB *sqlx.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	db, err := sqlx.Connect("postgres", cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &Postgres{DB: db}, nil
}

func (s *Postgres) Close() error {
	return s.DB.Close()
}

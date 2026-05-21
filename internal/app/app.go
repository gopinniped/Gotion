package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gopinniped/gotion/internal/config"
	"github.com/gopinniped/gotion/internal/domain/usecase"
	"github.com/gopinniped/gotion/internal/infra/postgres"
	"github.com/gopinniped/gotion/internal/infra/postgres/storage"
	"github.com/gopinniped/gotion/internal/transport/http/handler"
	"github.com/gopinniped/gotion/internal/transport/http/middleware"
	"github.com/gopinniped/gotion/pkg/token"
)

type App struct {
	db     *postgres.Postgres
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("init postgres: %w", err)
	}

	tokenMaker := token.NewMaker(cfg.JWTSecret)

	userStorage := storage.NewUserStorage(db.DB)
	userUseCase := usecase.NewUserUseCase(userStorage, tokenMaker)
	userHandler := handler.NewUserHandler(userUseCase)

	taskStorage := storage.NewTaskStorage(db.DB)
	taskUseCase := usecase.NewTaskUseCase(taskStorage)
	taskHandler := handler.NewTaskHandler(taskUseCase)

	router := NewRouter(userHandler, taskHandler, middleware.Auth(tokenMaker))

	server := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		db:     db,
		server: server,
	}, nil
}

func (a *App) Run() error {
	slog.Info("starting HTTP server", "addr", a.server.Addr)

	if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("http server: %w", err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	slog.Info("shutting down")

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	if err := a.db.Close(); err != nil {
		return fmt.Errorf("close db: %w", err)
	}

	return nil
}

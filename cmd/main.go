package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/gopinniped/gotion/internal/app"
	"github.com/gopinniped/gotion/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("init app: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := application.Run(); err != nil {
			log.Fatalf("run app: %v", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Stop(shutdownCtx); err != nil {
		log.Fatalf("stop app: %v", err)
	}
}

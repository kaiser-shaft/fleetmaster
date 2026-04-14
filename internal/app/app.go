package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaiser-shaft/fleetmaster/config"
	pgpool "github.com/kaiser-shaft/fleetmaster/pkg/postgres"
	redislib "github.com/kaiser-shaft/fleetmaster/pkg/redis"
)

func Run(ctx context.Context, c *config.Config) error {
	pgPool, err := pgpool.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	redisClient, err := redislib.New(c.Redis)
	if err != nil {
		return fmt.Errorf("redislib.New: %w", err)
	}

	// usecase
	// handler
	// httpserver

	slog.Info("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	slog.Info("App got signal to stop")

	pgPool.Close()
	redisClient.Close()

	slog.Info("App stopped!")

	return nil
}

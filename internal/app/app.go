package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaiser-shaft/fleetmaster/config"
)

func Run(ctx context.Context, c *config.Config) error {
	container := NewContainer(ctx, c)
	defer container.Close()

	_, err := container.HTTPServer()
	if err != nil {
		return fmt.Errorf("container.HTTPServer: %w", err)
	}

	slog.Info("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	slog.Info("App got signal to stop")

	return nil
}

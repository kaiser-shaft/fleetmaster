package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/kaiser-shaft/fleetmaster/config"
	"github.com/kaiser-shaft/fleetmaster/internal/app"
)

func main() {
	cfg := config.MustLoad()
	ctx := context.Background()

	if err := app.Run(ctx, cfg); err != nil {
		slog.Error("app.Run", slog.Any("error", err))
		os.Exit(1)
	}
}

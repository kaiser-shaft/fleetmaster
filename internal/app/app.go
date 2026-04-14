package app

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fleetmaster/config"
)

func Run(cfg *config.Config) {
	slog.Info("Starting server",
		slog.Int("port", cfg.HTTP.Port),
	)
	slog.Info("Random data", slog.String("UUID", uuid.New().String()))
}

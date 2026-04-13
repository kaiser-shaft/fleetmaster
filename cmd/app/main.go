package main

import (
	"log/slog"

	"github.com/kaiser-shaft/fleetmaster/internal/app"
)

func main() {
	slog.Info("Application started")

	app.Run()

	slog.Warn("Application stopped")
}

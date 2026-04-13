package main

import (
	"log/slog"

	"github.com/google/uuid"
)

func main() {
	slog.Info("Application started")

	slog.Info("Random data:", slog.String("UUID", uuid.New().String()))

	slog.Warn("Application stopped")
}

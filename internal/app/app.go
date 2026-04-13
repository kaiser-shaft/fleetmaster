package app

import (
	"log/slog"

	"github.com/google/uuid"
)

func Run() {
	slog.Info("Random data:", slog.String("UUID", uuid.New().String()))
}

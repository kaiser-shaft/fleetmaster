package main

import (
	"github.com/kaiser-shaft/fleetmaster/config"
	"github.com/kaiser-shaft/fleetmaster/internal/app"
)

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)
}

package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/kaiser-shaft/fleetmaster/pkg/httpserver"
	"github.com/kaiser-shaft/fleetmaster/pkg/postgres"
	"github.com/kaiser-shaft/fleetmaster/pkg/redis"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTP     httpserver.Config
	Postgres postgres.Config
	Redis    redis.Config
}

func New() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		if _, err := os.Stat(".env"); err == nil {
			configPath = ".env"
		}
	}

	var cfg Config
	if _, err := os.Stat(configPath); err == nil {
		if err := godotenv.Load(configPath); err != nil {
			return nil, fmt.Errorf("error loading %s file: %w", configPath, err)
		}
	}
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return &cfg, nil
}

func MustLoad() *Config {
	cfg, err := New()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

package redis

import (
	"log/slog"
	"net"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host string `envconfig:"REDIS_HOST" default:"localhost"`
	Port string `envconfig:"REDIS_PORT" default:"6379"`
}

type Client struct {
	*redis.Client
}

func New(c Config) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: net.JoinHostPort(c.Host, c.Port),
	})
	return &Client{Client: client}, nil
}

func (c *Client) Close() {
	err := c.Client.Close()
	if err != nil {
		slog.Error("redis.Close", slog.Any("error", err))
	}
}

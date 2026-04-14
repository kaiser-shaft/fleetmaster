package redis

type Config struct {
	Host string `envconfig:"REDIS_HOST" default:"localhost"`
	Port int    `envconfig:"REDIS_PORT" default:"6379"`
}

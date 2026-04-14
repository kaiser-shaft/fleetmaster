package httpserver

type Config struct {
	Port int `envconfig:"HTTP_PORT" default:"8080"`
}

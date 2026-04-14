package httpserver

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"
)

type Config struct {
	Port string `envconfig:"HTTP_PORT" default:"8080"`
}

type Server struct {
	server *http.Server
}

func New(handler http.Handler, c Config) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
		Addr:         net.JoinHostPort("", c.Port),
	}

	s := &Server{
		server: httpServer,
	}

	go s.start()

	return s
}

func (s *Server) start() {
	slog.Info("http server starting", slog.String("port", s.server.Addr))
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("httpserver.ListenAndServe", slog.Any("error", err))
	}
}

func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	slog.Info("http server shutting down...")

	if err := s.server.Shutdown(ctx); err != nil {
		slog.Warn("http server forced to shutdown", slog.Any("error", err))
	}
}

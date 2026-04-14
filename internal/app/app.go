package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaiser-shaft/fleetmaster/config"
	"github.com/kaiser-shaft/fleetmaster/internal/adapter/postgres"
	"github.com/kaiser-shaft/fleetmaster/internal/adapter/redis"
	v1 "github.com/kaiser-shaft/fleetmaster/internal/controller/http/v1"
	"github.com/kaiser-shaft/fleetmaster/internal/usecase"
	"github.com/kaiser-shaft/fleetmaster/pkg/httpserver"
	pgpool "github.com/kaiser-shaft/fleetmaster/pkg/postgres"
	redislib "github.com/kaiser-shaft/fleetmaster/pkg/redis"
	"net/http"
)

func Run(ctx context.Context, c *config.Config) error {
	pgPool, err := pgpool.New(ctx, c.Postgres)
	if err != nil {
		return fmt.Errorf("postgres.New: %w", err)
	}

	redisClient, err := redislib.New(c.Redis)
	if err != nil {
		return fmt.Errorf("redislib.New: %w", err)
	}

	// adapters
	userRepo := postgres.NewUserRepo(pgPool.Pool)
	vehicleRepo := postgres.NewVehicleRepo(pgPool.Pool)
	bookingRepo := postgres.NewBookingRepo(pgPool.Pool)
	cache := redis.NewCache(redisClient.Client)

	// usecase
	authUC := usecase.NewAuthUseCase(userRepo, cache)
	vehicleUC := usecase.NewVehicleUseCase(vehicleRepo)
	bookingUC := usecase.NewBookingUseCase(bookingRepo, vehicleRepo, userRepo, cache)

	// handlers
	authH := v1.NewAuthHandler(authUC)
	vehH := v1.NewVehicleHandler(vehicleUC)
	bookH := v1.NewBookingHandler(bookingUC)

	mux := http.NewServeMux()
	v1.NewRouter(mux, authUC, authH, vehH, bookH)

	// httpserver
	server := httpserver.New(mux, c.HTTP)

	slog.Info("App started!")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	slog.Info("App got signal to stop")

	server.Close()
	pgPool.Close()
	redisClient.Close()

	slog.Info("App stopped!")

	return nil
}

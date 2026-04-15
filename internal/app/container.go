package app

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/kaiser-shaft/fleetmaster/config"
	"github.com/kaiser-shaft/fleetmaster/internal/adapter/postgres"
	"github.com/kaiser-shaft/fleetmaster/internal/adapter/redis"
	v1 "github.com/kaiser-shaft/fleetmaster/internal/controller/http/v1"
	"github.com/kaiser-shaft/fleetmaster/internal/usecase"
	"github.com/kaiser-shaft/fleetmaster/pkg/httpserver"
	pgpool "github.com/kaiser-shaft/fleetmaster/pkg/postgres"
	redislib "github.com/kaiser-shaft/fleetmaster/pkg/redis"
)

type Container struct {
	ctx context.Context
	cfg *config.Config
	mu  sync.Mutex

	httpServer *httpserver.Server
	pgPool     *pgpool.Pool
	redis      *redislib.Client

	userRepo    *postgres.UserRepo
	vehicleRepo *postgres.VehicleRepo
	bookingRepo *postgres.BookingRepo
	cache       *redis.Cache

	authUC    *usecase.AuthUseCase
	vehicleUC *usecase.VehicleUseCase
	bookingUC *usecase.BookingUseCase
}

func NewContainer(ctx context.Context, c *config.Config) *Container {
	return &Container{
		ctx: ctx,
		cfg: c,
	}
}

func (c *Container) PGPool() (*pgpool.Pool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.pgPool == nil {
		pool, err := pgpool.New(c.ctx, c.cfg.Postgres)
		if err != nil {
			return nil, fmt.Errorf("container.PGPool: %w", err)
		}
		c.pgPool = pool
	}
	return c.pgPool, nil
}

func (c *Container) Redis() (*redislib.Client, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.redis == nil {
		client, err := redislib.New(c.cfg.Redis)
		if err != nil {
			return nil, fmt.Errorf("container.Redis: %w", err)
		}
		c.redis = client
	}
	return c.redis, nil
}

func (c *Container) UserRepo() (*postgres.UserRepo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.userRepo == nil {
		pool, err := c.getPGPool()
		if err != nil {
			return nil, err
		}
		c.userRepo = postgres.NewUserRepo(pool.Pool)
	}
	return c.userRepo, nil
}

func (c *Container) VehicleRepo() (*postgres.VehicleRepo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.vehicleRepo == nil {
		pool, err := c.getPGPool()
		if err != nil {
			return nil, err
		}
		c.vehicleRepo = postgres.NewVehicleRepo(pool.Pool)
	}
	return c.vehicleRepo, nil
}

func (c *Container) BookingRepo() (*postgres.BookingRepo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.bookingRepo == nil {
		pool, err := c.getPGPool()
		if err != nil {
			return nil, err
		}
		c.bookingRepo = postgres.NewBookingRepo(pool.Pool)
	}
	return c.bookingRepo, nil
}

func (c *Container) Cache() (*redis.Cache, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cache == nil {
		client, err := c.getRedis()
		if err != nil {
			return nil, err
		}
		c.cache = redis.NewCache(client.Client)
	}
	return c.cache, nil
}

func (c *Container) AuthUC() (*usecase.AuthUseCase, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.authUC == nil {
		ur, err := c.getUserRepo()
		if err != nil {
			return nil, err
		}
		ch, err := c.getCache()
		if err != nil {
			return nil, err
		}
		c.authUC = usecase.NewAuthUseCase(ur, ch)
	}
	return c.authUC, nil
}

func (c *Container) VehicleUC() (*usecase.VehicleUseCase, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.vehicleUC == nil {
		vr, err := c.getVehicleRepo()
		if err != nil {
			return nil, err
		}
		c.vehicleUC = usecase.NewVehicleUseCase(vr)
	}
	return c.vehicleUC, nil
}

func (c *Container) BookingUC() (*usecase.BookingUseCase, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.bookingUC == nil {
		br, err := c.getBookingRepo()
		if err != nil {
			return nil, err
		}
		vr, err := c.getVehicleRepo()
		if err != nil {
			return nil, err
		}
		ur, err := c.getUserRepo()
		if err != nil {
			return nil, err
		}
		ch, err := c.getCache()
		if err != nil {
			return nil, err
		}
		c.bookingUC = usecase.NewBookingUseCase(br, vr, ur, ch)
	}
	return c.bookingUC, nil
}

func (c *Container) HTTPServer() (*httpserver.Server, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.httpServer == nil {
		authUC, err := c.getAuthUC()
		if err != nil {
			return nil, err
		}
		vehicleUC, err := c.getVehicleUC()
		if err != nil {
			return nil, err
		}
		bookingUC, err := c.getBookingUC()
		if err != nil {
			return nil, err
		}

		authH := v1.NewAuthHandler(authUC)
		vehH := v1.NewVehicleHandler(vehicleUC)
		bookH := v1.NewBookingHandler(bookingUC)

		mux := http.NewServeMux()
		v1.NewRouter(mux, authUC, authH, vehH, bookH)

		c.httpServer = httpserver.New(mux, c.cfg.HTTP)
	}
	return c.httpServer, nil
}

// Вспомогательные приватные методы для использования внутри мьютекса
func (c *Container) getPGPool() (*pgpool.Pool, error) {
	if c.pgPool == nil {
		pool, err := pgpool.New(c.ctx, c.cfg.Postgres)
		if err != nil {
			return nil, err
		}
		c.pgPool = pool
	}
	return c.pgPool, nil
}

func (c *Container) getRedis() (*redislib.Client, error) {
	if c.redis == nil {
		client, err := redislib.New(c.cfg.Redis)
		if err != nil {
			return nil, err
		}
		c.redis = client
	}
	return c.redis, nil
}

func (c *Container) getUserRepo() (*postgres.UserRepo, error) {
	pool, err := c.getPGPool()
	if err != nil {
		return nil, err
	}
	if c.userRepo == nil {
		c.userRepo = postgres.NewUserRepo(pool.Pool)
	}
	return c.userRepo, nil
}

func (c *Container) getVehicleRepo() (*postgres.VehicleRepo, error) {
	pool, err := c.getPGPool()
	if err != nil {
		return nil, err
	}
	if c.vehicleRepo == nil {
		c.vehicleRepo = postgres.NewVehicleRepo(pool.Pool)
	}
	return c.vehicleRepo, nil
}

func (c *Container) getBookingRepo() (*postgres.BookingRepo, error) {
	pool, err := c.getPGPool()
	if err != nil {
		return nil, err
	}
	if c.bookingRepo == nil {
		c.bookingRepo = postgres.NewBookingRepo(pool.Pool)
	}
	return c.bookingRepo, nil
}

func (c *Container) getCache() (*redis.Cache, error) {
	client, err := c.getRedis()
	if err != nil {
		return nil, err
	}
	if c.cache == nil {
		c.cache = redis.NewCache(client.Client)
	}
	return c.cache, nil
}

func (c *Container) getAuthUC() (*usecase.AuthUseCase, error) {
	if c.authUC == nil {
		ur, err := c.getUserRepo()
		if err != nil {
			return nil, err
		}
		ch, err := c.getCache()
		if err != nil {
			return nil, err
		}
		c.authUC = usecase.NewAuthUseCase(ur, ch)
	}
	return c.authUC, nil
}

func (c *Container) getVehicleUC() (*usecase.VehicleUseCase, error) {
	if c.vehicleUC == nil {
		vr, err := c.getVehicleRepo()
		if err != nil {
			return nil, err
		}
		c.vehicleUC = usecase.NewVehicleUseCase(vr)
	}
	return c.vehicleUC, nil
}

func (c *Container) getBookingUC() (*usecase.BookingUseCase, error) {
	if c.bookingUC == nil {
		br, err := c.getBookingRepo()
		if err != nil {
			return nil, err
		}
		vr, err := c.getVehicleRepo()
		if err != nil {
			return nil, err
		}
		ur, err := c.getUserRepo()
		if err != nil {
			return nil, err
		}
		ch, err := c.getCache()
		if err != nil {
			return nil, err
		}
		c.bookingUC = usecase.NewBookingUseCase(br, vr, ur, ch)
	}
	return c.bookingUC, nil
}

func (c *Container) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.httpServer != nil {
		c.httpServer.Close()
	}
	if c.pgPool != nil {
		c.pgPool.Close()
	}
	if c.redis != nil {
		c.redis.Close()
	}
}

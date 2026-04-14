package repo

import (
	"context"

	"github.com/kaiser-shaft/fleetmaster/internal/entity"
)

type User interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}

type Vehicle interface {
	GetAll(ctx context.Context) ([]entity.Vehicle, error)
	GetByStatus(ctx context.Context, status entity.VehicleStatus) ([]entity.Vehicle, error)
	GetByID(ctx context.Context, id int64) (*entity.Vehicle, error)
	Update(ctx context.Context, vehicle *entity.Vehicle) error
}

type Booking interface {
	Create(ctx context.Context, booking *entity.Booking) error
	GetByID(ctx context.Context, id int64) (*entity.Booking, error)
	Update(ctx context.Context, booking *entity.Booking) error
	GetActiveByUserID(ctx context.Context, userID int64) (*entity.Booking, error)
}

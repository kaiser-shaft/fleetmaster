package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaiser-shaft/fleetmaster/internal/entity"
)

type BookingRepo struct {
	pool *pgxpool.Pool
}

func NewBookingRepo(pool *pgxpool.Pool) *BookingRepo {
	return &BookingRepo{pool: pool}
}

func (r *BookingRepo) Create(ctx context.Context, b *entity.Booking) error {
	err := r.pool.QueryRow(ctx, "INSERT INTO bookings (user_id, vehicle_id, start_time, end_time, purpose, status) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		b.UserID, b.VehicleID, b.StartTime, b.EndTime, b.Purpose, b.Status).Scan(&b.ID)
	if err != nil {
		return fmt.Errorf("BookingRepo.Create: %w", err)
	}
	return nil
}

func (r *BookingRepo) GetByID(ctx context.Context, id int64) (*entity.Booking, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, user_id, vehicle_id, start_time, end_time, purpose, status FROM bookings WHERE id = $1", id)
	var b entity.Booking
	err := row.Scan(&b.ID, &b.UserID, &b.VehicleID, &b.StartTime, &b.EndTime, &b.Purpose, &b.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("BookingRepo.GetByID: %w", err)
	}
	return &b, nil
}

func (r *BookingRepo) GetActiveByUserID(ctx context.Context, userID int64) (*entity.Booking, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, user_id, vehicle_id, start_time, end_time, purpose, status FROM bookings WHERE user_id = $1 AND status IN ('Pending', 'Active')", userID)
	var b entity.Booking
	err := row.Scan(&b.ID, &b.UserID, &b.VehicleID, &b.StartTime, &b.EndTime, &b.Purpose, &b.Status)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("BookingRepo.GetActiveByUserID: %w", err)
	}
	return &b, nil
}

func (r *BookingRepo) Update(ctx context.Context, b *entity.Booking) error {
	_, err := r.pool.Exec(ctx, "UPDATE bookings SET status = $1 WHERE id = $2", b.Status, b.ID)
	if err != nil {
		return fmt.Errorf("BookingRepo.Update: %w", err)
	}
	return nil
}

package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaiser-shaft/fleetmaster/internal/entity"
)

type VehicleRepo struct {
	pool *pgxpool.Pool
}

func NewVehicleRepo(pool *pgxpool.Pool) *VehicleRepo {
	return &VehicleRepo{pool: pool}
}

func (r *VehicleRepo) GetAll(ctx context.Context) ([]entity.Vehicle, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, brand, model, plate_number, status, mileage, last_service_mileage FROM vehicles")
	if err != nil {
		return nil, fmt.Errorf("VehicleRepo.GetAll: %w", err)
	}
	defer rows.Close()

	var vehicles []entity.Vehicle
	for rows.Next() {
		var v entity.Vehicle
		if err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.PlateNumber, &v.Status, &v.Mileage, &v.LastServiceMileage); err != nil {
			return nil, fmt.Errorf("VehicleRepo.GetAll.Scan: %w", err)
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (r *VehicleRepo) GetByStatus(ctx context.Context, status entity.VehicleStatus) ([]entity.Vehicle, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, brand, model, plate_number, status, mileage, last_service_mileage FROM vehicles WHERE status = $1", status)
	if err != nil {
		return nil, fmt.Errorf("VehicleRepo.GetByStatus: %w", err)
	}
	defer rows.Close()

	var vehicles []entity.Vehicle
	for rows.Next() {
		var v entity.Vehicle
		if err := rows.Scan(&v.ID, &v.Brand, &v.Model, &v.PlateNumber, &v.Status, &v.Mileage, &v.LastServiceMileage); err != nil {
			return nil, fmt.Errorf("VehicleRepo.GetByStatus.Scan: %w", err)
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (r *VehicleRepo) GetByID(ctx context.Context, id int64) (*entity.Vehicle, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, brand, model, plate_number, status, mileage, last_service_mileage FROM vehicles WHERE id = $1", id)
	var v entity.Vehicle
	err := row.Scan(&v.ID, &v.Brand, &v.Model, &v.PlateNumber, &v.Status, &v.Mileage, &v.LastServiceMileage)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("VehicleRepo.GetByID: %w", err)
	}
	return &v, nil
}

func (r *VehicleRepo) Update(ctx context.Context, v *entity.Vehicle) error {
	_, err := r.pool.Exec(ctx, "UPDATE vehicles SET status = $1, mileage = $2, last_service_mileage = $3 WHERE id = $4",
		v.Status, v.Mileage, v.LastServiceMileage, v.ID)
	if err != nil {
		return fmt.Errorf("VehicleRepo.Update: %w", err)
	}
	return nil
}

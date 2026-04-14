package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kaiser-shaft/fleetmaster/internal/entity"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, full_name, email, role, license_category FROM users WHERE email = $1", email)
	var u entity.User
	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.LicenseCategory)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepo.GetByEmail: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	row := r.pool.QueryRow(ctx, "SELECT id, full_name, email, role, license_category FROM users WHERE id = $1", id)
	var u entity.User
	err := row.Scan(&u.ID, &u.FullName, &u.Email, &u.Role, &u.LicenseCategory)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("UserRepo.GetByID: %w", err)
	}
	return &u, nil
}

func (r *UserRepo) Create(ctx context.Context, u *entity.User) error {
	err := r.pool.QueryRow(ctx, "INSERT INTO users (full_name, email, role, license_category) VALUES ($1, $2, $3, $4) RETURNING id",
		u.FullName, u.Email, u.Role, u.LicenseCategory).Scan(&u.ID)
	if err != nil {
		return fmt.Errorf("UserRepo.Create: %w", err)
	}
	return nil
}

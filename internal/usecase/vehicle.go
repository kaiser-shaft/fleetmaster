package usecase

import (
	"context"

	"github.com/kaiser-shaft/fleetmaster/internal/entity"
	"github.com/kaiser-shaft/fleetmaster/internal/usecase/repo"
)

type VehicleUseCase struct {
	repo repo.Vehicle
}

func NewVehicleUseCase(r repo.Vehicle) *VehicleUseCase {
	return &VehicleUseCase{repo: r}
}

func (uc *VehicleUseCase) GetAll(ctx context.Context) ([]entity.Vehicle, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *VehicleUseCase) GetByStatus(ctx context.Context, status entity.VehicleStatus) ([]entity.Vehicle, error) {
	return uc.repo.GetByStatus(ctx, status)
}

func (uc *VehicleUseCase) GetByID(ctx context.Context, id int64) (*entity.Vehicle, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *VehicleUseCase) SetRetired(ctx context.Context, id int64) error {
	v, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	v.Status = entity.StatusRetired
	return uc.repo.Update(ctx, v)
}

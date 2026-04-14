package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kaiser-shaft/fleetmaster/internal/entity"
	"github.com/kaiser-shaft/fleetmaster/internal/usecase/repo"
)

type Locker interface {
	AcquireLock(ctx context.Context, key string, expiration time.Duration) (bool, error)
	ReleaseLock(ctx context.Context, key string) error
}

type BookingUseCase struct {
	bookingRepo repo.Booking
	vehicleRepo repo.Vehicle
	userRepo    repo.User
	locker      Locker
}

func NewBookingUseCase(br repo.Booking, vr repo.Vehicle, ur repo.User, l Locker) *BookingUseCase {
	return &BookingUseCase{bookingRepo: br, vehicleRepo: vr, userRepo: ur, locker: l}
}

func (uc *BookingUseCase) Create(ctx context.Context, userID int64, vehicleID int64, startTime, endTime time.Time, purpose string) (*entity.Booking, error) {
	// Lock per user to prevent simultaneous bookings
	lockKey := fmt.Sprintf("booking_user:%d", userID)
	ok, err := uc.locker.AcquireLock(ctx, lockKey, 10*time.Second)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("booking in progress")
	}
	defer uc.locker.ReleaseLock(ctx, lockKey)

	// Rule 3.3: One active booking check
	active, err := uc.bookingRepo.GetActiveByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if active != nil {
		return nil, errors.New("user already has an active or pending booking")
	}

	// Rule 3.3: Vehicle status check
	vehicle, err := uc.vehicleRepo.GetByID(ctx, vehicleID)
	if err != nil {
		return nil, err
	}
	if vehicle == nil {
		return nil, errors.New("vehicle not found")
	}
	if vehicle.Status == entity.StatusMaintenance || vehicle.Status == entity.StatusRetired {
		return nil, fmt.Errorf("vehicle is in %s status and cannot be booked", vehicle.Status)
	}

	// Create booking
	booking := &entity.Booking{
		UserID:    userID,
		VehicleID: vehicleID,
		StartTime: startTime,
		EndTime:   endTime,
		Purpose:   purpose,
		Status:    entity.BookingPending,
	}

	// If start time is now or in the past, set to active and update vehicle status
	if !startTime.After(time.Now()) {
		booking.Status = entity.BookingActive
		vehicle.Status = entity.StatusInUse
		if err := uc.vehicleRepo.Update(ctx, vehicle); err != nil {
			return nil, err
		}
	}

	if err := uc.bookingRepo.Create(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (uc *BookingUseCase) Cancel(ctx context.Context, bookingID int64, userID int64) error {
	booking, err := uc.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if booking == nil {
		return errors.New("booking not found")
	}

	if booking.UserID != userID {
		return errors.New("cannot cancel someone else's booking")
	}

	if booking.Status != entity.BookingPending && booking.Status != entity.BookingActive {
		return errors.New("can only cancel pending or active bookings")
	}

	booking.Status = entity.BookingCancelled
	if err := uc.bookingRepo.Update(ctx, booking); err != nil {
		return err
	}

	vehicle, err := uc.vehicleRepo.GetByID(ctx, booking.VehicleID)
	if err == nil && vehicle != nil && vehicle.Status == entity.StatusInUse {
		vehicle.Status = entity.StatusAvailable
		return uc.vehicleRepo.Update(ctx, vehicle)
	}

	return nil
}

func (uc *BookingUseCase) Complete(ctx context.Context, bookingID int64, finalMileage int) error {
	booking, err := uc.bookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if booking == nil {
		return errors.New("booking not found")
	}
	if booking.Status != entity.BookingActive {
		return errors.New("can only complete active bookings")
	}

	vehicle, err := uc.vehicleRepo.GetByID(ctx, booking.VehicleID)
	if err != nil {
		return err
	}

	if finalMileage < vehicle.Mileage {
		return errors.New("final mileage cannot be less than current mileage")
	}
	vehicle.Mileage = finalMileage

	if vehicle.NeedsMaintenance() {
		vehicle.Status = entity.StatusMaintenance
	} else {
		vehicle.Status = entity.StatusAvailable
	}

	if err := uc.vehicleRepo.Update(ctx, vehicle); err != nil {
		return err
	}

	booking.Status = entity.BookingCompleted
	return uc.bookingRepo.Update(ctx, booking)
}

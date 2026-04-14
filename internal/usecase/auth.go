package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kaiser-shaft/fleetmaster/internal/entity"
	"github.com/kaiser-shaft/fleetmaster/internal/usecase/repo"
)

type SessionCache interface {
	SetSession(ctx context.Context, token string, userID int64, expiration time.Duration) error
	GetUserIDByToken(ctx context.Context, token string) (int64, error)
}

type AuthUseCase struct {
	userRepo     repo.User
	sessionCache SessionCache
}

func NewAuthUseCase(ur repo.User, sc SessionCache) *AuthUseCase {
	return &AuthUseCase{userRepo: ur, sessionCache: sc}
}

func (uc *AuthUseCase) Login(ctx context.Context, email string) (string, error) {
	user, err := uc.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", fmt.Errorf("user not found")
	}

	token := uuid.New().String()
	err = uc.sessionCache.SetSession(ctx, token, user.ID, 24*time.Hour)
	if err != nil {
		return "", fmt.Errorf("AuthUseCase.Login.SetSession: %w", err)
	}

	return token, nil
}

func (uc *AuthUseCase) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	return uc.sessionCache.GetUserIDByToken(ctx, token)
}

func (uc *AuthUseCase) GetUserByID(ctx context.Context, id int64) (*entity.User, error) {
	return uc.userRepo.GetByID(ctx, id)
}

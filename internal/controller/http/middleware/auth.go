package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/kaiser-shaft/fleetmaster/internal/entity"
	"github.com/kaiser-shaft/fleetmaster/pkg/render"
)

type SessionProvider interface {
	GetUserIDByToken(ctx context.Context, token string) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*entity.User, error)
}

type contextKey string

const (
	UserKey contextKey = "user"
)

func Auth(sp SessionProvider) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				render.Error(w, http.StatusUnauthorized, "missing authorization header", nil)
				return
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			userID, err := sp.GetUserIDByToken(r.Context(), token)
			if err != nil {
				render.Error(w, http.StatusUnauthorized, "invalid or expired token", err)
				return
			}

			user, err := sp.GetUserByID(r.Context(), userID)
			if err != nil || user == nil {
				render.Error(w, http.StatusUnauthorized, "user not found", err)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RoleRequired(role entity.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(UserKey).(*entity.User)
			if !ok || user.Role != role && user.Role != entity.RoleAdmin {
				render.Error(w, http.StatusForbidden, "insufficient permissions", nil)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func GetUser(r *http.Request) *entity.User {
	user, _ := r.Context().Value(UserKey).(*entity.User)
	return user
}

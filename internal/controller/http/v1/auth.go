package v1

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/kaiser-shaft/fleetmaster/pkg/render"
)

type AuthUC interface {
	Login(ctx context.Context, email string) (string, error)
}

type AuthHandler struct {
	uc AuthUC
}

func NewAuthHandler(uc AuthUC) *AuthHandler {
	return &AuthHandler{uc: uc}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	token, err := h.uc.Login(r.Context(), req.Email)
	if err != nil {
		render.Error(w, http.StatusUnauthorized, "login failed", err)
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"token": token})
}

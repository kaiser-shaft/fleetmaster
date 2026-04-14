package v1

import (
	"context"
	"net/http"
	"strconv"

	"github.com/kaiser-shaft/fleetmaster/internal/entity"
	"github.com/kaiser-shaft/fleetmaster/pkg/render"
)

type VehicleUC interface {
	GetAll(ctx context.Context) ([]entity.Vehicle, error)
	GetByStatus(ctx context.Context, status entity.VehicleStatus) ([]entity.Vehicle, error)
	GetByID(ctx context.Context, id int64) (*entity.Vehicle, error)
	SetRetired(ctx context.Context, id int64) error
}

type VehicleHandler struct {
	uc VehicleUC
}

func NewVehicleHandler(uc VehicleUC) *VehicleHandler {
	return &VehicleHandler{uc: uc}
}

func (h *VehicleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.uc.GetAll(r.Context())
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "failed to fetch vehicles", err)
		return
	}
	render.JSON(w, http.StatusOK, vehicles)
}

func (h *VehicleHandler) GetAvailable(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.uc.GetByStatus(r.Context(), entity.StatusAvailable)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "failed to fetch vehicles", err)
		return
	}
	render.JSON(w, http.StatusOK, vehicles)
}

func (h *VehicleHandler) GetMaintenance(w http.ResponseWriter, r *http.Request) {
	vehicles, err := h.uc.GetByStatus(r.Context(), entity.StatusMaintenance)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "failed to fetch vehicles", err)
		return
	}
	render.JSON(w, http.StatusOK, vehicles)
}

func (h *VehicleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	vehicle, err := h.uc.GetByID(r.Context(), id)
	if err != nil {
		render.Error(w, http.StatusInternalServerError, "failed to fetch vehicle", err)
		return
	}
	if vehicle == nil {
		render.Error(w, http.StatusNotFound, "vehicle not found", nil)
		return
	}
	render.JSON(w, http.StatusOK, vehicle)
}

func (h *VehicleHandler) Retire(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	if err := h.uc.SetRetired(r.Context(), id); err != nil {
		render.Error(w, http.StatusInternalServerError, "failed to retire vehicle", err)
		return
	}
	render.JSON(w, http.StatusOK, map[string]string{"status": "retired"})
}

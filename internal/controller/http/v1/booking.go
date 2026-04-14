package v1

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/kaiser-shaft/fleetmaster/internal/controller/http/middleware"
	"github.com/kaiser-shaft/fleetmaster/internal/entity"
	"github.com/kaiser-shaft/fleetmaster/pkg/render"
)

type BookingUC interface {
	Create(ctx context.Context, userID int64, vehicleID int64, startTime, endTime time.Time, purpose string) (*entity.Booking, error)
	Complete(ctx context.Context, bookingID int64, finalMileage int) error
	Cancel(ctx context.Context, bookingID int64, userID int64) error
}

type BookingHandler struct {
	uc BookingUC
}

func NewBookingHandler(uc BookingUC) *BookingHandler {
	return &BookingHandler{uc: uc}
}

func (h *BookingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VehicleID int64     `json:"vehicle_id"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
		Purpose   string    `json:"purpose"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	user := middleware.GetUser(r)
	booking, err := h.uc.Create(r.Context(), user.ID, req.VehicleID, req.StartTime, req.EndTime, req.Purpose)
	if err != nil {
		render.Error(w, http.StatusUnprocessableEntity, "could not create booking", err)
		return
	}

	render.JSON(w, http.StatusCreated, booking)
}

func (h *BookingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	user := middleware.GetUser(r)
	if err := h.uc.Cancel(r.Context(), id, user.ID); err != nil {
		render.Error(w, http.StatusUnprocessableEntity, "could not cancel booking", err)
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"status": "cancelled"})
}

func (h *BookingHandler) Complete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		render.Error(w, http.StatusBadRequest, "invalid ID", err)
		return
	}

	var req struct {
		Mileage int `json:"mileage"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		render.Error(w, http.StatusBadRequest, "invalid request body", err)
		return
	}

	if err := h.uc.Complete(r.Context(), id, req.Mileage); err != nil {
		render.Error(w, http.StatusUnprocessableEntity, "could not complete booking", err)
		return
	}

	render.JSON(w, http.StatusOK, map[string]string{"status": "completed"})
}

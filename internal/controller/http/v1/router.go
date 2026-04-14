package v1

import (
	"net/http"

	"github.com/kaiser-shaft/fleetmaster/internal/controller/http/middleware"
	"github.com/kaiser-shaft/fleetmaster/internal/entity"
)

func NewRouter(
	mux *http.ServeMux,
	sp middleware.SessionProvider,
	authH *AuthHandler,
	vehH *VehicleHandler,
	bookH *BookingHandler,
) {
	// Public routes
	mux.HandleFunc("POST /api/v1/auth/login", authH.Login)

	// Protected routes
	authMW := middleware.Auth(sp)
	adminMW := middleware.RoleRequired(entity.RoleAdmin)

	// Vehicles
	mux.Handle("GET /api/v1/vehicles", authMW(http.HandlerFunc(vehH.GetAll)))
	mux.Handle("GET /api/v1/vehicles/available", authMW(http.HandlerFunc(vehH.GetAvailable)))
	mux.Handle("GET /api/v1/vehicles/maintenance", authMW(adminMW(http.HandlerFunc(vehH.GetMaintenance))))
	mux.Handle("GET /api/v1/vehicles/{id}", authMW(http.HandlerFunc(vehH.GetByID)))
	mux.Handle("POST /api/v1/vehicles/{id}/retire", authMW(adminMW(http.HandlerFunc(vehH.Retire))))

	// Bookings
	mux.Handle("POST /api/v1/bookings", authMW(http.HandlerFunc(bookH.Create)))
	mux.Handle("POST /api/v1/bookings/{id}/complete", authMW(http.HandlerFunc(bookH.Complete)))
	mux.Handle("POST /api/v1/bookings/{id}/cancel", authMW(http.HandlerFunc(bookH.Cancel)))
}

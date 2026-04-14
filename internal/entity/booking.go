package entity

import "time"

type BookingStatus string

const (
	BookingPending   BookingStatus = "Pending"
	BookingActive    BookingStatus = "Active"
	BookingCompleted BookingStatus = "Completed"
	BookingCancelled BookingStatus = "Cancelled"
)

type Booking struct {
	ID        int64         `json:"id"`
	UserID    int64         `json:"user_id"`
	VehicleID int64         `json:"vehicle_id"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Purpose   string        `json:"purpose"`
	Status    BookingStatus `json:"status"`
}

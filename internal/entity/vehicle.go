package entity

type VehicleStatus string

const (
	StatusAvailable   VehicleStatus = "Available"
	StatusInUse       VehicleStatus = "In_Use"
	StatusMaintenance VehicleStatus = "Maintenance"
	StatusRetired     VehicleStatus = "Retired"
)

type Vehicle struct {
	ID                 int64         `json:"id"`
	Brand              string        `json:"brand"`
	Model              string        `json:"model"`
	PlateNumber        string        `json:"plate_number"`
	Status             VehicleStatus `json:"status"`
	Mileage            int           `json:"mileage"`
	LastServiceMileage int           `json:"last_service_mileage"`
}

func (v *Vehicle) NeedsMaintenance() bool {
	return v.Mileage-v.LastServiceMileage >= 10000
}

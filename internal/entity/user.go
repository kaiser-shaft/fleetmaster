package entity

type UserRole string

const (
	RoleAdmin   UserRole = "Admin"
	RoleDriver  UserRole = "Driver"
	RoleManager UserRole = "Manager"
)

type LicenseCategory string

const (
	LicenseA LicenseCategory = "A"
	LicenseB LicenseCategory = "B"
	LicenseC LicenseCategory = "C"
)

type User struct {
	ID              int64           `json:"id"`
	FullName        string          `json:"full_name"`
	Email           string          `json:"email"`
	Role            UserRole        `json:"role"`
	LicenseCategory LicenseCategory `json:"license_category"`
}
